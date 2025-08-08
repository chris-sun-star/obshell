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

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/oceanbase/obshell/agent/api/common"
	"github.com/oceanbase/obshell/agent/model/external"
	"github.com/oceanbase/obshell/agent/repository"
)

// @Summary Set Prometheus configuration
// @Description Set Prometheus configuration
// @Tags system
// @Accept json
// @Produce json
// @Param config body external.PrometheusConfig true "Prometheus configuration"
// @Success 200 {object} http.OcsAgentResponse
// @Failure 400 {object} http.OcsAgentResponse
// @Failure 500 {object} http.OcsAgentResponse
// @Router /api/v1/system/externals/prometheus [put]
func SetPrometheusConfig(c *gin.Context) {
	var cfg external.PrometheusConfig
	if err := c.BindJSON(&cfg); err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	repo, err := repository.NewExternalRepository()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	common.SendResponse(c, nil, repo.SavePrometheusConfig(&cfg))
}

// @Summary Get Prometheus configuration
// @Description Get Prometheus configuration
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} http.OcsAgentResponse{data=external.PrometheusConfig}
// @Failure 500 {object} http.OcsAgentResponse
// @Router /api/v1/system/externals/prometheus [get]
func GetPrometheusConfig(c *gin.Context) {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	cfg, err := repo.GetPrometheusConfig()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	common.SendResponse(c, cfg.Address, nil)
}

// @Summary Set Alertmanager configuration
// @Description Set Alertmanager configuration
// @Tags system
// @Accept json
// @Produce json
// @Param config body external.AlertmanagerConfig true "Alertmanager configuration"
// @Success 200 {object} http.OcsAgentResponse
// @Failure 400 {object} http.OcsAgentResponse
// @Failure 500 {object} http.OcsAgentResponse
// @Router /api/v1/system/externals/alertmanager [put]
func SetAlertmanagerConfig(c *gin.Context) {
	var cfg external.AlertmanagerConfig
	if err := c.BindJSON(&cfg); err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	repo, err := repository.NewExternalRepository()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	common.SendResponse(c, nil, repo.SaveAlertmanagerConfig(&cfg))
}

// @Summary Get Alertmanager configuration
// @Description Get Alertmanager configuration
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} http.OcsAgentResponse{data=external.AlertmanagerConfig}
// @Failure 500 {object} http.OcsAgentResponse
// @Router /api/v1/system/externals/alertmanager [get]
func GetAlertmanagerConfig(c *gin.Context) {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	cfg, err := repo.GetAlertmanagerConfig()
	if err != nil {
		common.SendResponse(c, nil, err)
		return
	}
	common.SendResponse(c, cfg.Address, nil)
}
