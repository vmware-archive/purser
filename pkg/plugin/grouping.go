/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	groups_v1 "github.com/vmware/purser/pkg/apis/groups/v1"
	groups "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
)

// GetGroupByName return group CRD by name.
func GetGroupByName(groupClient *groups.GroupClient, groupName string) *groups_v1.Group {
	group, err := groupClient.Get(groupName)
	if err != nil {
		log.Errorf("failed to get custom group by name %s, %v", groupName, err)
		return nil
	}
	return group
}

// PrintGroup displays the group information.
func PrintGroup(group *groups_v1.Group) {
	pitGroupMetrics := group.Spec.PITMetrics
	mtdGroupMetrics := group.Spec.MTDMetrics
	cost := group.Spec.MTDCost

	fmt.Printf("%-30s             %s\n", "Group Name:", group.Name)
	fmt.Println()

	if pitGroupMetrics != nil {
		fmt.Println("Point in Time Resource Stats:")
		fmt.Printf("             %-30s%.2f\n", "CPU Limit(vCPU):", pitGroupMetrics.CPULimit)
		fmt.Printf("             %-30s%.2f\n", "Memory Limit(GB):", pitGroupMetrics.MemoryLimit)
		fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU):", pitGroupMetrics.CPURequest)
		fmt.Printf("             %-30s%.2f\n", "Memory Request(GB):", pitGroupMetrics.MemoryRequest)
		fmt.Printf("             %-30s%.2f\n", "Storage Claimed(GB):", pitGroupMetrics.StorageClaim)
	}

	if mtdGroupMetrics != nil {
		fmt.Println()
		fmt.Printf("%-30s\n", "Month to Date Active Resource Stats:")
		fmt.Printf("             %-30s%.2f\n", "CPU Request(vCPU-hours):", mtdGroupMetrics.CPURequest)
		fmt.Printf("             %-30s%.2f\n", "Memory Request(GB-hours):", mtdGroupMetrics.MemoryRequest)
		fmt.Printf("             %-30s%.2f\n", "Storage Claimed(GB-hours):", mtdGroupMetrics.StorageClaim)
	}

	if cost != nil {
		fmt.Println()
		fmt.Printf("%-30s\n", "Month to Date Cost Stats:")
		fmt.Printf("             %-30s%.2f\n", "CPU Cost($):", cost.CPUCost)
		fmt.Printf("             %-30s%.2f\n", "Memory Cost($):", cost.MemoryCost)
		fmt.Printf("             %-30s%.2f\n", "Storage Cost($):", cost.StorageCost)
		fmt.Printf("             %-30s%.2f\n", "Total Cost($):", cost.TotalCost)
	}

	fmt.Println()
	fmt.Printf("Last updated timestamp(format RFC3339): %s", group.Spec.LastUpdated)
}
