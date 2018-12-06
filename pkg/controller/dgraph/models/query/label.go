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

// createLabelFilter will return a filter logic like
// (eq(key, "k1") AND eq(value, "v1")) OR (eq(key, "k1") AND eq(value, "v1")) OR (eq(key, "k1") AND eq(value, "v1"))
func createLabelFilter(labels map[string]string) string {
	seperator := " OR "
	var filter string
	isFirst := true
	for key, value := range labels {
		if !isFirst {
			filter += seperator
		} else {
			isFirst = false
		}
		filter += createSingleLabelFilter(key, value)
	}
	return filter
}

// createSingleLabelFilter takes key: k1, value: v1 and returns (eq(key, "k1") AND eq(value, "v1"))
func createSingleLabelFilter(key, value string) string {
	return `(eq(key, "` + key + `") AND eq(value, "` + value + `"))`
}
