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

package pool

import (
	"github.com/oceanbase/obshell/agent/engine/task"
	"github.com/oceanbase/obshell/agent/service/tenant"
)

const (
	TASK_NAME_DROP_RESOURCE_POOL = "Drop resource pools"

	PARAM_DROP_RESOURCE_POOL_LIST = "dropResourcePoolList"
)

var tenantService tenant.TenantService

func RegisterPoolTask() {
	task.RegisterTaskType(DropResourcePoolTask{})
}
