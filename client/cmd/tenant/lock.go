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
	"github.com/spf13/cobra"

	"github.com/oceanbase/obshell/agent/constant"
	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/lib/http"
	"github.com/oceanbase/obshell/client/command"
	clientconst "github.com/oceanbase/obshell/client/constant"
	"github.com/oceanbase/obshell/client/lib/stdio"
	"github.com/oceanbase/obshell/client/utils/api"
)

func newLockCmd() *cobra.Command {
	verbose := false
	lockCmd := command.NewCommand(&cobra.Command{
		Use:   CMD_LOCK,
		Short: "Lock a tenant.",
		Long:  "After lock tenant, the tenant will be read-only",
		RunE: command.WithErrorHandler(func(cmd *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.Occur(errors.ErrCliUsageError, "tenant name is required")
			}
			stdio.SetVerboseMode(verbose)
			return tenantLock(args[0])
		}),
		Example: `  obshell tenant lock t1`,
	})
	lockCmd.Annotations = map[string]string{clientconst.ANNOTATION_ARGS: "<tenant-name>"}
	lockCmd.VarsPs(&verbose, []string{clientconst.FLAG_VERBOSE, clientconst.FLAG_VERBOSE_SH}, false, "Activate verbose output", false)
	return lockCmd.Command
}

func tenantLock(name string) error {
	// Lock tenant
	stdio.StartLoadingf("lock tenant %s", name)
	if err := api.CallApiWithMethod(http.POST, constant.URI_TENANT_API_PREFIX+"/"+name+constant.URI_LOCK, nil, nil); err != nil {
		return err
	}
	stdio.LoadSuccessf("lock tenant %s", name)
	return nil
}
