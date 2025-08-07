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
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	alarmconstant "github.com/oceanbase/obshell/agent/executor/alarm/constant"
	"github.com/oceanbase/obshell/agent/model/external"
	"github.com/oceanbase/obshell/agent/repository"
	"github.com/pkg/errors"
)

func newClient(address string, auth *external.Auth) (*resty.Client, error) {
	client := resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second)).SetHostURL(address)
	if auth != nil && auth.Username != "" {
		client.SetBasicAuth(auth.Username, auth.Password)
	}
	return client, nil
}

func newAlertmanagerClient(cfg *external.AlertmanagerConfig) (*resty.Client, error) {
	return newClient(cfg.Address, cfg.Auth)
}

func newPrometheusClient(cfg *external.PrometheusConfig) (*resty.Client, error) {
	return newClient(cfg.Address, cfg.Auth)
}

func getAlertmanagerClientFromConfig() (*resty.Client, error) {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		return nil, errors.Wrap(err, "get external repository failed")
	}
	cfg, err := repo.GetAlertmanagerConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get alertmanager config failed")
	}
	if cfg == nil {
		return nil, errors.New("alertmanager config not found")
	}
	client, err := newAlertmanagerClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new alertmanager client failed")
	}
	return client, nil
}

func getPrometheusClientFromConfig() (*resty.Client, error) {
	repo, err := repository.NewExternalRepository()
	if err != nil {
		return nil, errors.Wrap(err, "get external repository failed")
	}
	cfg, err := repo.GetPrometheusConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get prometheus config failed")
	}
	if cfg == nil {
		return nil, errors.New("prometheus config not found")
	}
	client, err := newPrometheusClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new prometheus client failed")
	}
	return client, nil
}

func reloadAlertmanager() error {
	client, err := getAlertmanagerClientFromConfig()
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
	client, err := getPrometheusClientFromConfig()
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
