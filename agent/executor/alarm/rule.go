/*
Copyright (c) 2023 OceanBase
ob-operator is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/

package alarm

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/oceanbase/obshell/agent/executor/alarm/constant"
	"github.com/oceanbase/obshell/agent/repository"
	"github.com/oceanbase/obshell/model/alarm/rule"
	"github.com/pkg/errors"

	promv1 "github.com/prometheus/prometheus/web/api/v1"
	logger "github.com/sirupsen/logrus"
)

func GetRule(ctx context.Context, name string) (*rule.RuleResponse, error) {
	rules, err := ListRules(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "query rules from prometheus failed")
	}
	for _, rule := range rules {
		if rule.Name == name {
			return &rule, nil
		}
	}
	return nil, errors.New("rule not found")
}

func ListRules(ctx context.Context, filter *rule.RuleFilter) ([]rule.RuleResponse, error) {
	promRuleResponse := &rule.PromRuleResponse{}
	repo, err := repository.NewExternalRepository()
	if err != nil {
		return nil, errors.Wrap(err, "get external repository failed")
	}
	cfg, err := repo.GetPrometheusConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get prometheus config failed")
	}
	if cfg == nil {
		return nil, errors.New("prometheus config not found")
	}
	client, err := newPrometheusClient(cfg.Address, cfg.Auth.Username, cfg.Auth.Password)
	if err != nil {
		return nil, errors.Wrap(err, "new prometheus client failed")
	}
	resp, err := client.R().SetContext(ctx).SetQueryParam("type", "alert").SetHeader("content-type", "application/json").SetResult(promRuleResponse).Get(alarmconstant.RuleUrl)
	if err != nil {
		return nil, errors.Wrap(err, "query rules from prometheus failed")
	} else if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("query rules from prometheus got unexpected status: %d", resp.StatusCode())
	}
	logger.Debugf("Response from prometheus: %v", resp)
	filteredRules := make([]rule.RuleResponse, 0)
	for _, ruleGroup := range promRuleResponse.Data.RuleGroups {
		for _, promRule := range ruleGroup.Rules {
			encodedPromRule, err := json.Marshal(promRule)
			if err != nil {
				logger.Errorf("Got an error when encoding rule %v", promRule)
				continue
			}
			logger.Debugf("Process prometheus rule: %s", string(encodedPromRule))
			alertingRule := &promv1.AlertingRule{}
			err = json.Unmarshal(encodedPromRule, alertingRule)
			if err != nil {
				logger.Errorf("Got an error when decoding rule %v", promRule)
				continue
			}
			ruleResp := rule.NewRuleResponse(alertingRule)
			logger.Debugf("Parsed prometheus rule: %v", ruleResp)
			if filterRule(ruleResp, filter) {
				filteredRules = append(filteredRules, *ruleResp)
			}
		}
	}
	return filteredRules, nil
}

func filterRule(rule *rule.RuleResponse, filter *rule.RuleFilter) bool {
	matched := true
	if filter != nil {
		if filter.Keyword != "" {
			matched = matched && strings.Contains(rule.Name, filter.Keyword)
		}
		if filter.InstanceType != "" {
			matched = matched && (rule.InstanceType == filter.InstanceType)
		}
		if filter.Severity != "" {
			matched = matched && (rule.Severity == filter.Severity)
		}
	}
	return matched
}
