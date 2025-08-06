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

package repository

import (
	"encoding/json"

	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/model/external"
	obdb "github.com/oceanbase/obshell/agent/repository/db/oceanbase"
	obmodel "github.com/oceanbase/obshell/agent/repository/model/oceanbase"
	gorm "gorm.io/gorm"
)

const (
	PROMETHEUS_CONFIG_KEY = "prometheus_config"
	ALERTMANAGER_CONFIG_KEY = "alertmanager_config"
)

type ExternalRepository interface {
	SavePrometheusConfig(cfg *external.PrometheusConfig) error
	GetPrometheusConfig() (*external.PrometheusConfig, error)
	SaveAlertmanagerConfig(cfg *external.AlertmanagerConfig) error
	GetAlertmanagerConfig() (*external.AlertmanagerConfig, error)
}

type ExternalRepositoryImpl struct {
	db *gorm.DB
}

func NewExternalRepository() (*ExternalRepositoryImpl, error) {
	db, err := obdb.GetOceanbaseInstance()
	if err != nil {
		return nil, errors.Wrap(err, "get oceanbase instance failed")
	}
	return &ExternalRepositoryImpl{db: db}, nil
}

func (r *ExternalRepositoryImpl) SaveOcsConfig(name, value, info string) error {
	cfg := obmodel.OcsConfig{
		Name:  name,
		Value: value,
		Info:  info,
	}
	return r.db.Save(&cfg).Error
}

func (r *ExternalRepositoryImpl) GetOcsConfig(name string) (*obmodel.OcsConfig, error) {
	var cfg obmodel.OcsConfig
	err := r.db.Where("name = ?", name).First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "get ocs config failed")
	}
	return &cfg, nil
}

func (r *ExternalRepositoryImpl) SavePrometheusConfig(cfg *external.PrometheusConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, errors.ErrJsonMarshal.Code)
	}
	return r.SaveOcsConfig(PROMETHEUS_CONFIG_KEY, string(data), "Prometheus configuration")
}

func (r *ExternalRepositoryImpl) GetPrometheusConfig() (*external.PrometheusConfig, error) {
	ocsConfig, err := r.GetOcsConfig(PROMETHEUS_CONFIG_KEY)
	if err != nil {
		return nil, err
	}
	if ocsConfig == nil {
		return nil, nil
	}
	var cfg external.PrometheusConfig
	err = json.Unmarshal([]byte(ocsConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrJsonUnmarshal.Code)
	}
	return &cfg, nil
}

func (r *ExternalRepositoryImpl) SaveAlertmanagerConfig(cfg *external.AlertmanagerConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, errors.ErrJsonMarshal.Code)
	}
	return r.SaveOcsConfig(ALERTMANAGER_CONFIG_KEY, string(data), "Alertmanager configuration")
}

func (r *ExternalRepositoryImpl) GetAlertmanagerConfig() (*external.AlertmanagerConfig, error) {
	ocsConfig, err := r.GetOcsConfig(ALERTMANAGER_CONFIG_KEY)
	if err != nil {
		return nil, err
	}
	if ocsConfig == nil {
		return nil, nil
	}
	var cfg external.AlertmanagerConfig
	err = json.Unmarshal([]byte(ocsConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrJsonUnmarshal.Code)
	}
	return &cfg, nil
}