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

// Constants used in query parameters
const (
	AllPods  = ""
	Name     = "name"
	Orphan   = "orphan"
	View     = "view"
	Physical = "physical"
	Logical  = "logical"
	False    = "false"
)

// Children structure
type Children struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// Parent structure
type Parent struct {
	Name     string     `json:"name,omitempty"`
	Type     string     `json:"type,omitempty"`
	Children []Children `json:"children,omitempty"`
}

// JsonDataWrapper structure
type JsonDataWrapper struct {
	Data Parent `json:"data,omitempty"`
}
