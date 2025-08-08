/*
Copyright (c) 2023 OceanBase.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package alarm

import (
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
