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

package executor

import (
	"github.com/oceanbase/obshell/agent/executor/agent"
	"github.com/oceanbase/obshell/agent/executor/host"
	"github.com/oceanbase/obshell/agent/executor/ob"
	"github.com/oceanbase/obshell/agent/executor/obproxy"
	"github.com/oceanbase/obshell/agent/executor/pool"
	"github.com/oceanbase/obshell/agent/executor/recyclebin"
	"github.com/oceanbase/obshell/agent/executor/script"
	"github.com/oceanbase/obshell/agent/executor/task"
	"github.com/oceanbase/obshell/agent/executor/tenant"
	"github.com/oceanbase/obshell/agent/executor/unit"
	"github.com/oceanbase/obshell/agent/executor/zone"
)

func RegisterAllTask() {
	agent.RegisterAgentTask()
	ob.RegisterObInitTask()
	ob.RegisterObStartTask()
	ob.RegisterObStopTask()
	ob.RegisterObScaleOutTask()
	ob.RegisterObScaleInTask()
	ob.RegisterUpgradeTask()
	ob.RegisterBackupTask()
	ob.RegisterRestoreTask()
	obproxy.RegisterTaskType()
	pool.RegisterPoolTask()
	recyclebin.RegisterRecyclebinTask()
	script.RegisterScriptTask()
	task.RegisterTask()
	tenant.RegisterTenantTask()
	unit.RegisterUnitTask()
	zone.RegisterZoneTask()
	host.RegisterHostTask()
	RegisterMetricTask()
}
