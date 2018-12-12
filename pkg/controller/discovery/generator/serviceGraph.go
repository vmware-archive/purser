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

package generator

import (
	"strconv"

	"github.com/vmware/purser/pkg/controller/dgraph/models"
)


var (
	svcUniqueID int
	svcNodes    *[]Node
	svcEdges    *[]Edge
)

// GetGraphServiceNodes returns graph-nodes for service interactions
func GetGraphServiceNodes() []Node {
	return *svcNodes
}

// GetGraphServiceEdges returns graph-edges for service interactions
func GetGraphServiceEdges() []Edge {
	return *svcEdges
}

// GenerateServiceNodesAndEdges ...
func GenerateServiceNodesAndEdges(svcs []models.Service) {
	svcUniqueID = 0
	uniqueIDs := getServiceUniqueIDs(svcs)
	svcNodes := createServiceNodes(svcs, uniqueIDs)
	svcEdges := createServiceEdges(svcs, uniqueIDs)
	setServiceGraphNodes(svcNodes)
	setServiceGraphEdges(svcEdges)
}

func getServiceUniqueIDs(svcs []models.Service) map[string]int {
	uniqueIDs := make(map[string]int)
	for _, svc := range svcs {
		setServiceUniqueIDs(svc, uniqueIDs)
	}
	return uniqueIDs
}

func setServiceUniqueIDs(svc models.Service, uniqueIDs map[string]int) {
	if _, isPresent := uniqueIDs[svc.Name]; !isPresent {
		svcUniqueID++
		uniqueIDs[svc.Name] = svcUniqueID
	}
}

func createServiceNodes(svcs []models.Service, uniqueIDs map[string]int) []Node {
	nodes := []Node{}
	for _, svc := range svcs {
		newSvcNode := createNode("services", svc.Name, uniqueIDs[svc.Name], 0, []string{})
		nodes = append(nodes, newSvcNode)
	}
	return nodes
}

func createServiceEdges(svcs []models.Service, uniqueIDs map[string]int) []Edge {
	edges := []Edge{}
	for _, svc := range svcs {
		srcID := uniqueIDs[svc.Name]
		for _, dstSvc := range svc.Services {
			destID := uniqueIDs[dstSvc.Name]
			edges = append(edges, createEdge(srcID, destID, 0))
		}
	}
	return edges
}

func setServiceGraphNodes(nodes []Node) {
	svcNodes = &nodes
}

func setServiceGraphEdges(edges []Edge) {
	svcEdges = &edges
}
