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

package zone

import (
	"strings"

	"github.com/oceanbase/obshell/agent/constant"
	"github.com/oceanbase/obshell/agent/errors"
	"github.com/oceanbase/obshell/agent/repository/model/oceanbase"
	tenantservice "github.com/oceanbase/obshell/agent/service/tenant"
	"github.com/oceanbase/obshell/param"
	"github.com/oceanbase/obshell/utils"
)

func CheckPrimaryZoneAndLocality(primaryZone string, locality map[string]string) error {
	// Get first priority zones.
	firstPriorityZones := make([]string, 0)
	if primaryZone == constant.PRIMARY_ZONE_RANDOM {
		for zone := range locality {
			firstPriorityZones = append(firstPriorityZones, zone)
		}
	} else {
		firstPriorityZones = strings.Split(strings.Split(primaryZone, ";")[0], ",")
	}

	// Build zone -> region map
	zonesWithRegion, err := obclusterService.GetAllZonesWithRegion()
	if err != nil {
		return err
	}
	zoneToRegionMap := make(map[string]string, 0)
	for _, z := range zonesWithRegion {
		zoneToRegionMap[z.Zone] = z.Region
	}

	// Check whether first priority zones are in the same region.
	var firstPriorityRegion string
	for _, zone := range firstPriorityZones {
		if firstPriorityRegion == "" {
			firstPriorityRegion = zoneToRegionMap[zone]
		} else if firstPriorityRegion != zoneToRegionMap[zone] {
			return errors.New("Tenant primary zone could not span regions.")
		}
	}

	// Check whether the locality has multi-region.
	firstPriorityRegion = zoneToRegionMap[firstPriorityZones[0]]
	hasMultiRegion := false
	for zone := range locality {
		if zoneToRegionMap[zone] != firstPriorityRegion {
			hasMultiRegion = true
			break
		}
	}
	// If there is only one region, no need to check the number of full replicas.
	if !hasMultiRegion {
		return nil
	}

	// The first priority region should have more than 1 full replica when locality has multi-region.
	fullReplicaNum := 0
	for zone, replicaType := range locality {
		if zoneToRegionMap[zone] == firstPriorityRegion {
			arr := strings.Split(replicaType, "{")
			if arr[0] == constant.REPLICA_TYPE_FULL || arr[0] == "F" || arr[0] == "" {
				fullReplicaNum++
			}
		}
	}
	if fullReplicaNum < 2 {
		return errors.Errorf("The region %v where the first priority of tenant zone is located needs to have at least 2 F replicas. In fact, there are only %d full replicas.", firstPriorityRegion, fullReplicaNum)
	}

	return nil
}

func GetFirstPriorityPrimaryZone(locality string, primaryZone string) ([]string, error) {
	replicaInfoMap, err := tenantservice.ParseLocalityToReplicaInfoMap(locality)
	if err != nil {
		return nil, nil
	}
	var firstPriorityZones []string
	if primaryZone == constant.PRIMARY_ZONE_RANDOM {
		for zone, replicaType := range replicaInfoMap {
			if replicaType == constant.REPLICA_TYPE_FULL {
				// Only full replica zone can be first priority primary zone.
				firstPriorityZones = append(firstPriorityZones, zone)
			}
		}
	} else {
		zoneArrayList := strings.Split(primaryZone, ";")
		for _, zoneArray := range zoneArrayList {
			firstPriorityZones = make([]string, 0)
			zoneList := strings.Split(zoneArray, ",")
			for _, zone := range zoneList {
				if replicaType, ok := replicaInfoMap[zone]; ok && replicaType == constant.REPLICA_TYPE_FULL {
					// Only full replica zone can be first priority primary zone.
					firstPriorityZones = append(firstPriorityZones, zone)
				}
			}
			if len(firstPriorityZones) > 0 {
				break
			}
			// if there is no full replica zone in this zone array, then continue to check next zone array.
		}
	}
	return firstPriorityZones, nil
}

func isFirstPriorityPrimaryZoneChangedWhenAlterParimaryZone(tenant *oceanbase.DbaObTenant, targetPrimaryZone string) (bool, error) {
	prevFirstPriorityZones, err := GetFirstPriorityPrimaryZone(tenant.Locality, tenant.PrimaryZone)
	if err != nil {
		return false, err
	}
	newFirstPriorityZones, err := GetFirstPriorityPrimaryZone(tenant.Locality, targetPrimaryZone)
	if err != nil {
		return false, err
	}
	return !utils.SliceEqual(prevFirstPriorityZones, newFirstPriorityZones), nil
}

func isFirstPriorityPrimaryZoneChangedWhenAlterLocality(tenant *oceanbase.DbaObTenant, targetLocality string) (bool, error) {
	prevFirstPriorityZones, err := GetFirstPriorityPrimaryZone(tenant.Locality, tenant.PrimaryZone)
	if err != nil {
		return false, err
	}
	newFirstPriorityZones, err := GetFirstPriorityPrimaryZone(targetLocality, tenant.PrimaryZone)
	if err != nil {
		return false, err
	}
	return !utils.SliceEqual(prevFirstPriorityZones, newFirstPriorityZones), nil
}

func CheckFirstPriorityPrimaryZoneChangedWhenAlterPrimaryZone(tenant *oceanbase.DbaObTenant, targetPrimaryZone string) error {
	changed, err := isFirstPriorityPrimaryZoneChangedWhenAlterParimaryZone(tenant, targetPrimaryZone)
	if err != nil {
		return err
	}
	if changed {
		if enableRebalance, err := tenantService.GetTenantParameter(tenant.TenantID, constant.PARAMETER_ENABLE_REBALANCE); err != nil {
			return err
		} else if enableRebalance == nil {
			return errors.New("Get enable_rebalance failed.")
		} else if enableRebalance.Value != "True" {
			return errors.New("Change first priority zone of primary zone is not allowed when tenant 'enable_rebalance' is disabled")
		}
	}
	return nil
}

func CheckFirstPriorityPrimaryZoneChangedWhenAlterLocality(tenant *oceanbase.DbaObTenant, targetLocality string) error {
	changed, err := isFirstPriorityPrimaryZoneChangedWhenAlterLocality(tenant, targetLocality)
	if err != nil {
		return err
	}
	if changed {
		if enableRebalance, err := tenantService.GetTenantParameter(tenant.TenantID, constant.PARAMETER_ENABLE_REBALANCE); err != nil {
			return err
		} else if enableRebalance == nil {
			return errors.New("Get enable_rebalance failed.")
		} else if enableRebalance.Value != "True" {
			return errors.New("Change first priority zone of primary zone is not allowed when tenant 'enable_rebalance' is disabled")
		}
	}
	return nil
}

func CheckPrimaryZone(primaryZone string, zoneList []string) error {
	if primaryZone == constant.PRIMARY_ZONE_RANDOM {
		return nil
	}
	zonesSemicolonSeparated := strings.Split(primaryZone, ";")
	exsitZones := make([]string, 0)
	for _, zones := range zonesSemicolonSeparated {
		zonesCommaSeparated := strings.Split(zones, ",")
		for _, zone := range zonesCommaSeparated {
			if !utils.ContainsString(zoneList, zone) {
				return errors.Errorf("Zone '%s' is not in zone_list.", zone)
			} else if utils.ContainsString(exsitZones, zone) {
				return errors.Errorf("Zone '%s' is repeated in primary_zone.", zone)
			} else {
				exsitZones = append(exsitZones, zone)
			}
		}
	}
	return nil
}

func CheckAtLeastOnePaxosReplica(zoneList []param.ZoneParam) error {
	for _, zone := range zoneList {
		if zone.ReplicaType == constant.REPLICA_TYPE_FULL {
			return nil
		}
	}
	return errors.New("At least one zone should be FULL replica.")
}

func CheckZoneParams(zoneList []param.ZoneParam) error {
	if len(zoneList) == 0 {
		return errors.New("zone_list is empty")
	}

	if err := StaticCheckForZoneParams(zoneList); err != nil {
		return err
	}

	for _, zone := range zoneList {
		// Check whether the zone exists
		if exist, err := obclusterService.IsZoneExist(zone.Name); err != nil {
			return err
		} else if !exist {
			return errors.Errorf("Zone '%s' is not exist.", zone.Name)
		}

		// Check unit config if exsits.
		if exist, err := unitService.IsUnitConfigExist(zone.UnitConfigName); err != nil {
			return err
		} else if !exist {
			return errors.Errorf("Unit config '%s' is not exist.", zone.UnitConfigName)
		}

		servers, err := obclusterService.GetServerByZone(zone.Name)
		if err != nil {
			return err
		}
		if len(servers) < zone.UnitNum {
			return errors.Errorf("The number of servers in zone '%s' is %d, less than the number of units %d.", zone.Name, len(servers), zone.UnitNum)
		}
	}
	return nil
}

func StaticCheckForZoneParams(zoneList []param.ZoneParam) error {
	unitNum := 0
	existZones := make([]string, 0)
	for _, zone := range zoneList {
		if utils.ContainsString(existZones, zone.Name) {
			return errors.Errorf("Zone '%s' is repeated.", zone.Name)
		}
		existZones = append(existZones, zone.Name)

		if zone.UnitConfigName == "" {
			return errors.New("unit_config_name should not be empty.")
		}

		// Check replica type.
		if err := CheckReplicaType(zone.ReplicaType); err != nil {
			return err
		}

		// Check unit num.
		if zone.UnitNum <= 0 {
			return errors.New("unit_num should be positive.")
		}

		if zone.UnitNum != unitNum && unitNum != 0 {
			return errors.New("unit_num should be same in all zones.")
		}
		unitNum = zone.UnitNum
	}
	return nil
}

func CheckReplicaType(localityType string) error {
	if localityType != constant.REPLICA_TYPE_FULL && localityType != constant.REPLICA_TYPE_READONLY && localityType != "" {
		return errors.New("ReplicaType should be 'FULL' or 'READONLY'")
	}
	return nil
}
