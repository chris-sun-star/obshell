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

package tenant

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/oceanbase/obshell/agent/config"
	"github.com/oceanbase/obshell/agent/constant"
	"github.com/oceanbase/obshell/agent/engine/task"
	"github.com/oceanbase/obshell/agent/lib/http"
	ocsagentlog "github.com/oceanbase/obshell/agent/log"
	"github.com/oceanbase/obshell/client/command"
	clientconst "github.com/oceanbase/obshell/client/constant"
	"github.com/oceanbase/obshell/client/global"
	"github.com/oceanbase/obshell/client/lib/stdio"
	"github.com/oceanbase/obshell/client/utils/api"
)

func newPurgeCmd() *cobra.Command {
	opts := global.DropFlags{}
	purgeCmd := command.NewCommand(&cobra.Command{
		Use:   CMD_PURGE,
		Short: "Purge a tenant in recyclebin.",
		Long:  "Resource will be free later",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			if len(args) <= 0 {
				stdio.Error("tenant or object name is required")
				cmd.SilenceUsage = false
				return errors.New("tenant or object name is required")
			}
			cmd.SilenceUsage = true
			ocsagentlog.InitLogger(config.DefaultClientLoggerConifg())
			stdio.SetVerboseMode(opts.Verbose)
			stdio.SetSkipConfirmMode(opts.SkipConfirm)
			if err := tenantPurge(args[0]); err != nil {
				stdio.LoadFailedWithoutMsg()
				stdio.Error(err.Error())
				return err
			}
			return nil
		},
		Example: `  obshell recyclebin tenant purge t1
  obshell recyclebin tenant purge '__recycle_$_1_1720679549921648'`,
	})

	purgeCmd.Annotations = map[string]string{clientconst.ANNOTATION_ARGS: "<tenant-name|object-name>"}
	purgeCmd.Flags().SortFlags = false
	purgeCmd.VarsPs(&opts.Verbose, []string{clientconst.FLAG_VERBOSE, clientconst.FLAG_VERBOSE_SH}, false, "Activate verbose output", false)
	purgeCmd.VarsPs(&opts.SkipConfirm, []string{clientconst.FLAG_SKIP_CONFIRM, clientconst.FLAG_SKIP_CONFIRM_SH}, false, "Skip the confirmation of purge tenant operation", false)
	return purgeCmd.Command
}

func tenantPurge(name string) error {
	pass, err := stdio.Confirmf("Please confirm if you need to purge tenant %s", name)
	if err != nil {
		return errors.New("ask for confirmation failed")
	}
	if !pass {
		return nil
	}
	var dag task.DagDetailDTO
	// Drop tenant
	if err := api.CallApiWithMethod(http.DELETE, constant.URI_API_V1+constant.URI_RECYCLEBIN_GROUP+constant.URI_TENANT_GROUP+"/"+name, nil, &dag); err != nil {
		return err
	}
	if dag.GenericDTO == nil {
		stdio.Printf("No such tenant '%s' in recyclebin", name)
		return nil
	}
	if err := api.NewDagHandler(&dag).PrintDagStage(); err != nil {
		return err
	}
	return nil
}
