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

package models

import (
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
	"golang.org/x/crypto/bcrypt"
)

// Login constants
const (
	DefaultUsername = "admin"
	DefaultPassword = "purser!123"
	DefaultLoginXID = "purser-login-xid"
	IsLogin         = "isLogin"
)

// Login structure
type Login struct {
	dgraph.ID
	IsLogin  bool   `json:"isLogin,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// StoreLogin ...
func StoreLogin() {
	uid := dgraph.GetUID(DefaultLoginXID, IsLogin)
	if uid == "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), bcrypt.MinCost)
		if err != nil {
			logrus.Errorf("error while hashing login information")
		}
		login := Login{
			ID:       dgraph.ID{Xid: DefaultLoginXID},
			IsLogin:  true,
			Username: DefaultUsername,
			Password: string(hashedPassword),
		}
		_, err = dgraph.MutateNode(login, dgraph.CREATE)
		if err != nil {
			logrus.Errorf("error while storing login information")
		}
	}
}

// UpdateLogin ...
func UpdateLogin(username, oldPassword, newPassword string) bool {
	if query.CheckLogin(username, oldPassword) {
		login, err := query.GetHashedPassword(username)
		if err != nil {
			logrus.Error(err)
			return false
		}
		err = hashAndUpdatePassword(&login, newPassword)
		if err == nil {
			return true
		}
		logrus.Error(err)
	}
	return false
}

func hashAndUpdatePassword(login *Login, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}
	login.Password = string(hashedPassword)
	_, err = dgraph.MutateNode(login, dgraph.UPDATE)
	return err
}
