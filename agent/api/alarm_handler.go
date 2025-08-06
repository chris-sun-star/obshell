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
)

// ListAlerts godoc
// @ID ListAlerts
// @Summary List all alerts
// @Description List all alerts
// @Tags alarm
// @Accept json
// @Produce json
// @Param filter body alert.AlertFilter true "alert filter"
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
// @Param filter body silence.SilencerFilter true "silencer filter"
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
// @Param silencer body silence.SilencerParam true "silencer param"
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
// @Param filter body rule.RuleFilter true "rule filter"
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
