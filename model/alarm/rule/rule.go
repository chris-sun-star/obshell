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

package rule

import (
	promv1 "github.com/prometheus/prometheus/web/api/v1"
	alarm "github.com/oceanbase/obshell/model/alarm"
)

type Rule struct {
	Name        string `json:"name" binding:"required"`
	InstanceType string `json:"instanceType" binding:"required"`
	Severity    alarm.Severity `json:"severity" binding:"required"`
	// Add other fields as needed based on usage in ob-operator/internal/dashboard/business/alarm/rule.go
}

type RuleResponse struct {
	Name        string `json:"name"`
	InstanceType string `json:"instanceType"`
	Severity    alarm.Severity `json:"severity"`
	// Add other fields as needed
}

func NewRuleResponse(alertingRule *promv1.AlertingRule) *RuleResponse {
	// Implement conversion from promv1.AlertingRule to RuleResponse
	return &RuleResponse{
		Name: alertingRule.Name,
		// Populate other fields based on alertingRule
	}
}

type RuleFilter struct {
	Keyword      string `json:"keyword"`
	InstanceType string `json:"instanceType"`
	Severity     string `json:"severity"`
}

type PromRuleResponse struct {
	Data struct {
		RuleGroups []struct {
			Rules []interface{} `json:"rules"` // Use interface{} as the actual type is promv1.Rule
		} `json:"groups"`
	} `json:"data"`
}
