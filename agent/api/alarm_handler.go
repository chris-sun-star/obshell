/*
Copyright (c) 2024 OceanBase.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/oceanbase/obshell/agent/executor/alarm"
	"github.com/oceanbase/obshell/model/alarm/alert"
	"github.com/oceanbase/obshell/model/alarm/rule"
	"github.com/oceanbase/obshell/model/alarm/silence"
)

// ListAlerts godoc
// @ID ListAlerts
// @Summary List all alerts
// @Description List all alerts
// @Tags alarm
// @Accept json
// @Produce json
// @Param filter body alert.AlertFilter false "alert filter"
// @Success 200 {object} http.OcsAgentResponse{data=[]alert.Alert}
// @Router /api/v1/alarm/alert/alerts [post]
func ListAlerts(ctx *gin.Context) ([]alert.Alert, error) {
	filter := &alert.AlertFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, err
	}
	return alarm.ListAlerts(ctx, filter)
}

// ListSilencers godoc
// @ID ListSilencers
// @Summary List all silencers
// @Description List all silencers
// @Tags alarm
// @Accept json
// @Produce json
// @Param filter body silence.SilencerFilter false "silencer filter"
// @Success 200 {object} http.OcsAgentResponse{data=[]silence.SilencerResponse}
// @Router /api/v1/alarm/silence/silencers [post]
func ListSilencers(ctx *gin.Context) ([]silence.SilencerResponse, error) {
	filter := &silence.SilencerFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, err
	}
	return alarm.ListSilencers(ctx, filter)
}

// GetSilencer godoc
// @ID GetSilencer
// @Summary Get a silencer
// @Description Get a silencer by id
// @Tags alarm
// @Accept json
// @Produce json
// @Param id path string true "silencer id"
// @Success 200 {object} http.OcsAgentResponse{data=silence.SilencerResponse}
// @Router /api/v1/alarm/silence/silencers/{id} [get]
func GetSilencer(ctx *gin.Context) (*silence.SilencerResponse, error) {
	id := ctx.Param("id")
	return alarm.GetSilencer(ctx, id)
}

// CreateOrUpdateSilencer godoc
// @ID CreateOrUpdateSilencer
// @Summary Create or update a silencer
// @Description Create or update a silencer
// @Tags alarm
// @Accept json
// @Produce json
// @Param silencer body silence.SilencerParam true "silencer"
// @Success 200 {object} http.OcsAgentResponse{data=silence.SilencerResponse}
// @Router /api/v1/alarm/silence/silencers [put]
func CreateOrUpdateSilencer(ctx *gin.Context) (*silence.SilencerResponse, error) {
	param := &silence.SilencerParam{}
	err := ctx.Bind(param)
	if err != nil {
		return nil, err
	}
	return alarm.CreateOrUpdateSilencer(ctx, param)
}

// DeleteSilencer godoc
// @ID DeleteSilencer
// @Summary Delete a silencer
// @Description Delete a silencer by id
// @Tags alarm
// @Accept json
// @Produce json
// @Param id path string true "silencer id"
// @Success 204
// @Router /api/v1/alarm/silence/silencers/{id} [delete]
func DeleteSilencer(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return nil, alarm.DeleteSilencer(ctx, id)
}

// ListRules godoc
// @ID ListRules
// @Summary List all rules
// @Description List all rules
// @Tags alarm
// @Accept json
// @Produce json
// @Param filter body rule.RuleFilter false "rule filter"
// @Success 200 {object} http.OcsAgentResponse{data=[]rule.RuleResponse}
// @Router /api/v1/alarm/rule/rules [post]
func ListRules(ctx *gin.Context) ([]rule.RuleResponse, error) {
	filter := &rule.RuleFilter{}
	err := ctx.Bind(filter)
	if err != nil {
		return nil, err
	}
	return alarm.ListRules(ctx, filter)
}

// GetRule godoc
// @ID GetRule
// @Summary Get a rule
// @Description Get a rule by name
// @Tags alarm
// @Accept json
// @Produce json
// @Param name path string true "rule name"
// @Success 200 {object} http.OcsAgentResponse{data=rule.RuleResponse}
// @Router /api/v1/alarm/rule/rules/{name} [get]
func GetRule(ctx *gin.Context) (*rule.RuleResponse, error) {
	name := ctx.Param("name")
	return alarm.GetRule(ctx, name)
}
