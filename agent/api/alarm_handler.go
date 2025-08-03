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

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/oceanbase/obshell/agent/executor/alarm"
	"github.com/oceanbase/obshell/model/alarm/alert"
	"github.com/oceanbase/obshell/model/alarm/rule"
	"github.com/oceanbase/obshell/model/alarm/silence"
	"github.com/pkg/errors"
)

// ListAlerts
func ListAlerts(ctx *gin.Context) ([]alert.Alert, error) {
	filter := &alert.AlertFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	return alarm.ListAlerts(ctx, filter)
}

// ListSilencers
func ListSilencers(ctx *gin.Context) ([]silence.SilencerResponse, error) {
	filter := &silence.SilencerFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	return alarm.ListSilencers(ctx, filter)
}

// GetSilencer
func GetSilencer(ctx *gin.Context) (*silence.SilencerResponse, error) {
	id := ctx.Param("id")
	return alarm.GetSilencer(ctx, id)
}

// CreateOrUpdateSilencer
func CreateOrUpdateSilencer(ctx *gin.Context) (*silence.SilencerResponse, error) {
	param := &silence.SilencerParam{}
	err := ctx.Bind(param)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	return alarm.CreateOrUpdateSilencer(ctx, param)
}

// DeleteSilencer
func DeleteSilencer(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return nil, alarm.DeleteSilencer(ctx, id)
}

// ListRules
func ListRules(ctx *gin.Context) ([]rule.RuleResponse, error) {
	filter := &rule.RuleFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	return alarm.ListRules(ctx, filter)
}

// GetRule
func GetRule(ctx *gin.Context) (*rule.RuleResponse, error) {
	name := ctx.Param("name")
	return alarm.GetRule(ctx, name)
}
