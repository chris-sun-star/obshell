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
package ob

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/oceanbase/obshell/agent/constant"
	"github.com/oceanbase/obshell/agent/lib/parse"
	"github.com/oceanbase/obshell/param"
)

func GetClusterUnitSpecLimit() *param.ClusterUnitConfigLimit {
	minMemory, err := obclusterService.GetParameterByName(constant.PARAMETER_MIN_FULL_RESOURCE_POOL_MEMORY)
	if err != nil {
		log.Warnf("get %s failed, err: %s", constant.PARAMETER_MIN_FULL_RESOURCE_POOL_MEMORY, err.Error())
		return nil
	}
	var unitSpecLimit param.ClusterUnitConfigLimit
	if minMemory != nil {
		minMemoryValue, err := strconv.Atoi(minMemory.Value)
		if err != nil {
			log.Warnf("convert %s to int failed, err: %s", minMemory.Value, err.Error())
			return nil
		}
		unitSpecLimit.MinMemory = float64(minMemoryValue) / parse.GB
	}
	unitSpecLimit.MinCpu = 1
	return &unitSpecLimit
}
