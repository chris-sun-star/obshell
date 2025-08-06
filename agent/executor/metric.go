/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package executor

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oceanbase/obshell/agent/bindata"
	"github.com/oceanbase/obshell/agent/constant"
	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/repository"
	"github.com/oceanbase/obshell/model/common"
	model "github.com/oceanbase/obshell/model/metric"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	PROMETHEUS_ADDRESS      = "http://127.0.0.1:9090"
	METRIC_RANGE_QUERY_URL  = "/api/v1/query_range"
	DEFAULT_TIMEOUT         = 30
	METRIC_CONFIG_FILE_ENUS = "agent/assets/metric/metrics-en_US.yaml"
	METRIC_CONFIG_FILE_ZHCN = "agent/assets/metric/metrics-zh_CN.yaml"
	METRIC_EXPR_CONFIG_FILE = "agent/assets/metric/metric_expr.yaml"
	KEY_INTERVAL            = "@INTERVAL"
	KEY_LABELS              = "@LABELS"
	KEY_GROUP_LABELS        = "@GBLABELS"
)

var metricExprConfig map[string]string

func init() {
	metricExprConfig = make(map[string]string)
	metricExprConfigContent, err := bindata.Asset(METRIC_EXPR_CONFIG_FILE)
	if err != nil {
		logrus.WithError(err).Error("load metric expr config failed")
	}
	err = yaml.Unmarshal(metricExprConfigContent, &metricExprConfig)
	if err != nil {
		logrus.WithError(err).Error("parse metric expr config data failed")
	}
}

func ListMetricClasses(scope, language string) ([]model.MetricClass, error) {
	metricClasses := make([]model.MetricClass, 0)
	configFile := METRIC_CONFIG_FILE_ENUS
	switch language {
	case constant.LANGUAGE_EN_US:
		configFile = METRIC_CONFIG_FILE_ENUS
	case constant.LANGUAGE_ZH_CN:
		configFile = METRIC_CONFIG_FILE_ZHCN
	default:
		logrus.Infof("Not supported language %s, return default", language)
	}

	metricConfigContent, err := bindata.Asset(configFile)
	if err != nil {
		return metricClasses, err
	}
	metricConfigMap := make(map[string][]model.MetricClass)
	// TODO: Do not unmarshal the file every time, cache the result
	err = yaml.Unmarshal(metricConfigContent, &metricConfigMap)
	if err != nil {
		return metricClasses, err
	}
	logrus.Debugf("metric configs: %v", metricConfigMap)
	metricClasses, found := metricConfigMap[scope]
	if !found {
		err = errors.Errorf("metric config for scope %s not found", scope)
	}
	return metricClasses, err
}

func replaceQueryVariables(exprTemplate string, labels []common.KVPair, groupLabels []string, step int64) string {
	labelStrParts := make([]string, 0, len(labels))
	for _, label := range labels {
		labelStrParts = append(labelStrParts, fmt.Sprintf("%s=\"%s\"", label.Key, label.Value))
	}
	labelStr := strings.Join(labelStrParts, ",")
	groupLabelStr := strings.Join(groupLabels, ",")
	replacer := strings.NewReplacer(KEY_INTERVAL, fmt.Sprintf("%ss", strconv.FormatInt(step, 10)), KEY_LABELS, labelStr, KEY_GROUP_LABELS, groupLabelStr)
	return replacer.Replace(exprTemplate)
}

func extractMetricData(name string, resp *model.PrometheusQueryRangeResponse) []model.MetricData {
	metricDatas := make([]model.MetricData, 0)
	for _, result := range resp.Data.Result {
		values := make([]model.MetricValue, 0)

		labels := make([]common.KVPair, 0, len(result.Metric))
		for k, v := range result.Metric {
			labels = append(labels, common.KVPair{Key: k, Value: v})
		}

		metric := model.Metric{
			Name:   name,
			Labels: labels,
		}
		lastValid := math.NaN()
		invalidTimestamps := make([]float64, 0)
		// one loop to handle invalid timestamps interpolation
		for _, value := range result.Values {
			t := value[0].(float64)
			v, err := strconv.ParseFloat(value[1].(string), 64)
			if err != nil {
				logrus.Warnf("failed to parse value %v", value)
				invalidTimestamps = append(invalidTimestamps, t)
			} else if math.IsNaN(v) {
				logrus.Debugf("value at timestamp %f is NaN", t)
				invalidTimestamps = append(invalidTimestamps, t)
			} else {
				// if there are invalid timestamps, interpolate them
				if len(invalidTimestamps) > 0 {
					var interpolated float64
					if math.IsNaN(lastValid) {
						interpolated = v
					} else {
						interpolated = (lastValid + v) / 2
					}
					// interpolate invalid slots with last valid value
					for _, it := range invalidTimestamps {
						values = append(values, model.MetricValue{
							Timestamp: it,
							Value:     interpolated,
						})
					}
					invalidTimestamps = invalidTimestamps[:0]
				}
				values = append(values, model.MetricValue{
					Timestamp: t,
					Value:     v,
				})
				lastValid = v
			}
		}
		if math.IsNaN(lastValid) {
			lastValid = 0.0
		}
		for _, it := range invalidTimestamps {
			values = append(values, model.MetricValue{
				Timestamp: it,
				Value:     lastValid,
			})
		}
		metricDatas = append(metricDatas, model.MetricData{
			Metric: metric,
			Values: values,
		})
	}
	return metricDatas
}

func QueryMetricData(queryParam *model.MetricQuery) []model.MetricData {
	client := resty.New().SetTimeout(time.Duration(DEFAULT_TIMEOUT * time.Second))
	repo, err := NewExternalRepository()
	if err != nil {
		logrus.WithError(err).Error("get external repository failed")
		return nil
	}
	cfg, err := repo.GetPrometheusConfig()
	if err != nil {
		logrus.WithError(err).Error("get prometheus config failed")
		return nil
	}
	if cfg == nil {
		logrus.Error("prometheus config not found")
		return nil
	}
	if cfg.User != "" {
		client.SetBasicAuth(cfg.User, cfg.Password)
	}

	metricDatas := make([]model.MetricData, 0, len(queryParam.Metrics))
	wg := sync.WaitGroup{}
	metricDataCh := make(chan []model.MetricData, len(queryParam.Metrics))
	for _, m := range queryParam.Metrics {
		exprTemplate, found := metricExprConfig[m]
		if found {
			wg.Add(1)
			go func(m string, ch chan []model.MetricData) {
				defer wg.Done()
				expr := replaceQueryVariables(exprTemplate, queryParam.Labels, queryParam.GroupLabels, queryParam.QueryRange.Step)
				logrus.Infof("Query with expr: %s, range: %v", expr, queryParam.QueryRange)
				queryRangeResp := &model.PrometheusQueryRangeResponse{}
				resp, err := client.R().SetQueryParams(map[string]string{
					"start": strconv.FormatFloat(queryParam.QueryRange.StartTimestamp, 'f', 3, 64),
					"end":   strconv.FormatFloat(queryParam.QueryRange.EndTimestamp, 'f', 3, 64),
					"step":  strconv.FormatInt(queryParam.QueryRange.Step, 10),
					"query": expr,
				}).SetHeader("content-type", "application/json").
					SetResult(queryRangeResp).
					Get(fmt.Sprintf("%s%s", cfg.URL, METRIC_RANGE_QUERY_URL))
				if err != nil {
					logrus.Errorf("Query expression expr got error: %v", err)
				} else if resp.StatusCode() == http.StatusOK {
					ch <- extractMetricData(m, queryRangeResp)
				}
			}(m, metricDataCh)
		} else {
			logrus.Warnf("Metric %s expression not found", m)
		}
	}
	wg.Wait()
	close(metricDataCh)
	for metricDataArray := range metricDataCh {
		metricDatas = append(metricDatas, metricDataArray...)
	}
	return metricDatas
}
