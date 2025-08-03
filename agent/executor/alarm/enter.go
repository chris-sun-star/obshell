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
	"github.com/oceanbase/obshell/agent/engine/task"
)

func RegisterAlarmTask() {
	task.RegisterTaskType(ListAlertsTask{})
	task.RegisterTaskType(ListSilencersTask{})
	task.RegisterTaskType(GetSilencerTask{})
	task.RegisterTaskType(CreateOrUpdateSilencerTask{})
	task.RegisterTaskType(DeleteSilencerTask{})
	task.RegisterTaskType(ListRulesTask{})
	task.RegisterTaskType(GetRuleTask{})
}

type ListAlertsTask struct { task.Task }

func (t *ListAlertsTask) Execute() error { return nil }

type ListSilencersTask struct { task.Task }

func (t *ListSilencersTask) Execute() error { return nil }

type GetSilencerTask struct { task.Task }

func (t *GetSilencerTask) Execute() error { return nil }

type CreateOrUpdateSilencerTask struct { task.Task }

func (t *CreateOrUpdateSilencerTask) Execute() error { return nil }

type DeleteSilencerTask struct { task.Task }

func (t *DeleteSilencerTask) Execute() error { return nil }

type ListRulesTask struct { task.Task }

func (t *ListRulesTask) Execute() error { return nil }

type GetRuleTask struct { task.Task }

func (t *GetRuleTask) Execute() error { return nil }
