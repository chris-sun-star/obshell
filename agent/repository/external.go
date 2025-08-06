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

	"github.com/oceanbase/obshell/agent/config"
	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/repository/db/sqlite"
	model "github.com/oceanbase/obshell/agent/repository/model/sqlite"
	gorm "gorm.io/gorm"
)

const (
	PROMETHEUS_CONFIG_KEY = "prometheus_config"
	ALERTMANAGER_CONFIG_KEY = "alertmanager_config"
)

type ExternalRepository interface {
	SavePrometheusConfig(cfg *config.PrometheusConfig) error
	GetPrometheusConfig() (*config.PrometheusConfig, error)
	SaveAlertmanagerConfig(cfg *config.AlertmanagerConfig) error
	GetAlertmanagerConfig() (*config.AlertmanagerConfig, error)
}

type ExternalRepositoryImpl struct {
	db *gorm.DB
}

func NewExternalRepository() (*ExternalRepositoryImpl, error) {
	db, err := sqlite.GetSqliteInstance()
	if err != nil {
		return nil, errors.Wrap(err, "get sqlite instance failed")
	}
	return &ExternalRepositoryImpl{db: db}, nil
}

func (r *ExternalRepositoryImpl) SaveObConfig(name, value, info string) error {
	cfg := model.ObConfig{
		Name:  name,
		Value: value,
		Info:  info,
	}
	return r.db.Save(&cfg).Error
}

func (r *ExternalRepositoryImpl) GetObConfig(name string) (*model.ObConfig, error) {
	var cfg model.ObConfig
	err := r.db.Where("name = ?", name).First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "get ob config failed")
	}
	return &cfg, nil
}

func (r *ExternalRepositoryImpl) SavePrometheusConfig(cfg *config.PrometheusConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, errors.ErrJsonMarshal.Code)
	}
	return r.SaveObConfig(PROMETHEUS_CONFIG_KEY, string(data), "Prometheus configuration")
}

func (r *ExternalRepositoryImpl) GetPrometheusConfig() (*config.PrometheusConfig, error) {
	obConfig, err := r.GetObConfig(PROMETHEUS_CONFIG_KEY)
	if err != nil {
		return nil, err
	}
	if obConfig == nil {
		return nil, nil
	}
	var cfg config.PrometheusConfig
	err = json.Unmarshal([]byte(obConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrJsonUnmarshal.Code)
	}
	return &cfg, nil
}

func (r *ExternalRepositoryImpl) SaveAlertmanagerConfig(cfg *config.AlertmanagerConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, errors.ErrJsonMarshal.Code)
	}
	return r.SaveObConfig(ALERTMANAGER_CONFIG_KEY, string(data), "Alertmanager configuration")
}

func (r *ExternalRepositoryImpl) GetAlertmanagerConfig() (*config.AlertmanagerConfig, error) {
	obConfig, err := r.GetObConfig(ALERTMANAGER_CONFIG_KEY)
	if err != nil {
		return nil, err
	}
	if obConfig == nil {
		return nil, nil
	}
	var cfg config.AlertmanagerConfig
	err = json.Unmarshal([]byte(obConfig.Value), &cfg)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrJsonUnmarshal.Code)
	}
	return &cfg, nil
}
