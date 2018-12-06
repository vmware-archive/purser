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

import (
	"github.com/vmware/purser/test/utils"
	"testing"
)

// TestCreateSingleLabelFilter ...
func TestCreateSingleLabelFilter(t *testing.T) {
	got := createSingleLabelFilter("k1", "v1")
	expected := `(eq(key, "k1") AND eq(value, "v1"))`
	utils.Equals(t, expected, got)
}

// TestCreateLabelFilter ...
func TestCreateLabelFilter(t *testing.T) {
	labels := make(map[string]string)
	labels["k1"] = "v1"
	got := createLabelFilter(labels)
	expected := `(eq(key, "k1") AND eq(value, "v1"))`
	utils.Equals(t, expected, got)

	labels["k2"] = "v2"
	got2 := createLabelFilter(labels)
	expected1 := `(eq(key, "k2") AND eq(value, "v2")) OR (eq(key, "k1") AND eq(value, "v1"))`
	expected2 := `(eq(key, "k1") AND eq(value, "v1")) OR (eq(key, "k2") AND eq(value, "v2"))`
	utils.Assert(t, (got2 == expected1) || (got2 == expected2), "label filter didn't match")
}
