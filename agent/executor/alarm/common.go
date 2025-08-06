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

	"github.com/go-resty/resty/v2"
	"github.com/oceanbase/obshell/agent/executor/alarm/constant"
	metricconst "github.com/oceanbase/obshell/agent/executor/metric/constant"
	"github.com/oceanbase/obshell/agent/repository"
	"github.com/pkg/errors"
)

func newClient(url, user, password string) (*resty.Client, error) {
	client := resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second)).SetHostURL(url)
	if user != "" {
		client.SetBasicAuth(user, password)
	}
	return client, nil
}

func newAlertmanagerClient(url, user, password string) (*resty.Client, error) {
	return newClient(url, user, password)
}

func newPrometheusClient(url, user, password string) (*resty.Client, error) {
	return newClient(url, user, password)
}

func reloadAlertmanager() error {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		return errors.Wrap(err, "get external repository failed")
	}
	cfg, err := repo.GetAlertmanagerConfig()
	if err != nil {
		return errors.Wrap(err, "get alertmanager config failed")
	}
	if cfg == nil {
		return errors.New("alertmanager config not found")
	}
	client, err := newAlertmanagerClient(cfg.URL, cfg.User, cfg.Password)
	if err != nil {
		return errors.Wrap(err, "new alertmanager client failed")
	}
	resp, err := client.R().SetHeader("content-type", "application/json").Post(alarmconstant.AlertmanagerReloadUrl)
	if err != nil {
		return errors.Wrap(err, "reload alertmanager failed")
	} else if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("reload alertmanager got unexpected status: %d", resp.StatusCode())
	}
	return nil
}

func reloadPrometheus() error {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		return errors.Wrap(err, "get external repository failed")
	}
	cfg, err := repo.GetPrometheusConfig()
	if err != nil {
		return errors.Wrap(err, "get prometheus config failed")
	}
	if cfg == nil {
		return errors.New("prometheus config not found")
	}
	client, err := newPrometheusClient(cfg.URL, cfg.User, cfg.Password)
	if err != nil {
		return errors.Wrap(err, "new prometheus client failed")
	}
	resp, err := client.R().SetHeader("content-type", "application/json").Post(alarmconstant.PrometheusReloadUrl)
	if err != nil {
		return errors.Wrap(err, "reload prometheus failed")
	} else if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("reload prometheus got unexpected status: %d", resp.StatusCode())
	}
	return nil
}
