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
	"fmt"
	"net/http"
	"sync"
	"time"

	alarmconstant "github.com/oceanbase/obshell/agent/executor/alarm/constant"
	metricconst "github.com/oceanbase/obshell/agent/executor/metric/constant"
	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

var restyClient *resty.Client
var restyOnce sync.Once

func getClient() *resty.Client {
	restyOnce.Do(func() {
		restyClient = resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second))
	})
	return restyClient
}

func reloadAlertmanager() error {
	client := resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second))
	resp, err := client.R().SetHeader("content-type", "application/json").Post(fmt.Sprintf("%s%s", alarmconstant.AlertManagerAddress, alarmconstant.AlertmanagerReloadUrl))
	if err != nil {
		return errors.Wrap(err, errors.ErrExternal, "Reload alertmanager failed")
	} else if resp.StatusCode() != http.StatusOK {
		return errors.Newf(errors.ErrExternal, "Reload alertmanager got unexpected status: %d", resp.StatusCode())
	}
	return nil
}

func reloadPrometheus() error {
	client := resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second))
	resp, err := client.R().SetHeader("content-type", "application/json").Post(fmt.Sprintf("%s%s", metricconst.PrometheusAddress, alarmconstant.PrometheusReloadUrl))
	if err != nil {
		return errors.Wrap(err, errors.ErrExternal, "Reload prometheus failed")
	} else if resp.StatusCode() != http.StatusOK {
		return errors.Newf(errors.ErrExternal, "Reload prometheus got unexpected status: %d", resp.StatusCode())
	}
	return nil
}
