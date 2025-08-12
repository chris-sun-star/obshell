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

package config

import (
	"encoding/json"

	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/model/external"
	configservice "github.com/oceanbase/obshell/agent/service/config"
)

const (
	PROMETHEUS_CONFIG_KEY   = "prometheus_config"
	ALERTMANAGER_CONFIG_KEY = "alertmanager_config"
)

func SavePrometheusConfig(cfg *external.PrometheusConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Occur(errors.ErrJsonMarshal, err.Error())
	}
	return configservice.SaveOcsConfig(PROMETHEUS_CONFIG_KEY, string(data), "Prometheus configuration")
}

func GetPrometheusConfig() (*external.PrometheusConfig, error) {
	ocsConfig, err := configservice.GetOcsConfig(PROMETHEUS_CONFIG_KEY)
	if err != nil {
		return nil, errors.WrapRetain(errors.ErrConfigGetFailed, err, PROMETHEUS_CONFIG_KEY, err.Error())
	}
	if ocsConfig == nil {
		return nil, errors.Occurf(errors.ErrConfigNotFound, PROMETHEUS_CONFIG_KEY)
	}
	var cfg external.PrometheusConfig
	err = json.Unmarshal([]byte(ocsConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Occur(errors.ErrJsonUnmarshal, err.Error())
	}
	return &cfg, nil
}

func SaveAlertmanagerConfig(cfg *external.AlertmanagerConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Occur(errors.ErrJsonMarshal, err.Error())
	}
	return configservice.SaveOcsConfig(ALERTMANAGER_CONFIG_KEY, string(data), "Alertmanager configuration")
}

func GetAlertmanagerConfig() (*external.AlertmanagerConfig, error) {
	ocsConfig, err := configservice.GetOcsConfig(ALERTMANAGER_CONFIG_KEY)
	if err != nil {
		return nil, errors.WrapRetain(errors.ErrConfigGetFailed, err, ALERTMANAGER_CONFIG_KEY, err.Error())
	}
	if ocsConfig == nil {
		return nil, nil
	}
	var cfg external.AlertmanagerConfig
	err = json.Unmarshal([]byte(ocsConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Occur(errors.ErrJsonUnmarshal, err.Error())
	}
	return &cfg, nil
}
