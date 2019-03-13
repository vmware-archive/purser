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

package query

const (
	testSecondsSinceMonthStart = "1.45"
	testPodUIDList             = "0x3e283, 0x3e288"
	testPodName                = "pod-purser-dgraph-0"
	testDaemonsetName          = "daemonset-purser"
	testResourceName           = "resource-purser"
	testPodUID                 = "0x3e283"
	testPodXID                 = "purser:pod-purser-dgraph-0"

	testHierarchy            = "hierarchy"
	testMetrics              = "metrics"
	testRetrieveAllGroups    = "retrieveAllGroups"
	testRetrieveGroupMetrics = "retrieveGroupMetrics"
	testRetrieveSubscribers  = "retrieveSubscribers"
	testLabelFilterPods      = "labelFilterPods"
	testAlivePods            = "alivePods"
	testPodInteractions      = "podInteractions"
	testPodPrices            = "podPrices"
	testCapacity             = "capacityAllocation"
	testWrongQuery           = "wrongQuery"
	testCPUPrice             = 0.24
	testMemoryPrice          = 0.1
)
