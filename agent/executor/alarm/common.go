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

package alarm

import (
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/oceanbase/obshell/agent/errors"
	alarmconstant "github.com/oceanbase/obshell/agent/executor/alarm/constant"
	configexecutor "github.com/oceanbase/obshell/agent/executor/config"
	"github.com/oceanbase/obshell/agent/model/external"
)

func newClient(address string, auth *external.Auth) *resty.Client {
	client := resty.New().SetTimeout(time.Duration(alarmconstant.DefaultAlarmQueryTimeout * time.Second)).SetHostURL(address)
	if auth != nil && auth.Username != "" {
		client.SetBasicAuth(auth.Username, auth.Password)
	}
	return client
}

func newAlertmanagerClient(cfg *external.AlertmanagerConfig) *resty.Client {
	return newClient(cfg.Address, cfg.Auth)
}

func newPrometheusClient(cfg *external.PrometheusConfig) *resty.Client {
	return newClient(cfg.Address, cfg.Auth)
}

func getAlertmanagerClientFromConfig() (*resty.Client, error) {
	cfg, err := configexecutor.GetAlertmanagerConfig()
	if err != nil {
		return nil, errors.WrapRetain(errors.ErrConfigGetFailed, err, configexecutor.ALERTMANAGER_CONFIG_KEY, err.Error())
	}
	if cfg == nil {
		return nil, errors.Occur(errors.ErrConfigNotFound, configexecutor.ALERTMANAGER_CONFIG_KEY)
	}
	client := newAlertmanagerClient(cfg)
	return client, nil
}

func getPrometheusClientFromConfig() (*resty.Client, error) {
	cfg, err := configexecutor.GetPrometheusConfig()
	if err != nil {
		return nil, errors.WrapRetain(errors.ErrConfigGetFailed, err, configexecutor.PROMETHEUS_CONFIG_KEY, err.Error())
	}
	if cfg == nil {
		return nil, errors.Occur(errors.ErrConfigNotFound, configexecutor.ALERTMANAGER_CONFIG_KEY)
	}
	client := newPrometheusClient(cfg)
	return client, nil
}
