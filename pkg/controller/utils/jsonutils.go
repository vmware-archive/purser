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

package utils

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// JSONMarshal marshal object and returns in byte. If there is an error then it return nil.
func JSONMarshal(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error(err)
	}
	return bytes
}

// GetJSONResponse retrieves json response and converts it to target object.
// Returns error if any failure is encountered.
func GetJSONResponse(client *http.Client, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	return json.NewDecoder(resp.Body).Decode(target)
}

func closeResponse(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Errorf("unable to close response body. Reason: %v", err)
	}
}
