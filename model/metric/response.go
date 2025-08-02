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

package metric

type MetricValue struct {
	Timestamp float64 `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Metric struct {
	Name   string   `json:"name"`
	Labels []KVPair `json:"labels"`
}

type MetricData struct {
	Metric Metric        `json:"metric"`
	Values []MetricValue `json:"values"`
}

type MetricInfo struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type MetricGroup struct {
	Name    string       `json:"name"`
	Metrics []MetricInfo `json:"metrics"`
}

type MetricClass struct {
	Name   string        `json:"name"`
	Groups []MetricGroup `json:"groups"`
}
