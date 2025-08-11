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
	"github.com/oceanbase/obshell/agent/constant"
)

func InitAlarmRoutes(parentGroup *gin.RouterGroup, isLocalRoute bool) {
	alarm := parentGroup.Group(constant.URI_ALARM_GROUP)

	if !isLocalRoute {
		alarm.Use(common.Verify())
	}

	// alerts
	alert := alarm.Group(constant.URI_ALERT_GROUP)
	alert.POST(constant.URI_ALERTS, ListAlerts)

	// silencers
	silence := alarm.Group(constant.URI_SILENCE_GROUP)
	silence.POST(constant.URI_SILENCERS, ListSilencers)
	silence.GET(constant.URI_SILENCERS+constant.URI_PATH_PARAM_ID, GetSilencer)
	silence.PUT(constant.URI_SILENCERS, CreateOrUpdateSilencer)
	silence.DELETE(constant.URI_SILENCERS+constant.URI_PATH_PARAM_ID, DeleteSilencer)

	// rules
	rule := alarm.Group(constant.URI_RULE_GROUP)
	rule.POST(constant.URI_RULES, ListRules)
	rule.GET(constant.URI_RULES+constant.URI_PATH_PARAM_NAME, GetRule)
}
