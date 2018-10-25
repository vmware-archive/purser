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
	"os"
	"os/user"

	log "github.com/Sirupsen/logrus"
)

// OpenFile handles opening file in Read/Write mode, creating and appending to it as needed.
func OpenFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("failed to open file %s, %v", filename, err)
	}
	return f
}

// GetUsrHomeDir returns the current user's Home Directory
func GetUsrHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Errorf("failed to fetch current user %v", err)
	}
	return usr.HomeDir
}
