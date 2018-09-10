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

package main

// These are possible actions for resources
const (
	Get = "get"
	Set = "set"
)

// These are kubernetes components
const (
	Label     = "label"
	Pod       = "pod"
	Node      = "node"
	Namespace = "namespace"
	Group     = "group"
)

// These are utilisation metrics
const (
	Cost      = "cost"
	Resources = "resources"
)
