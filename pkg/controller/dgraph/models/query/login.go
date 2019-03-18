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
	"github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/dgraph"

	"golang.org/x/crypto/bcrypt"
)

// Login structure
type Login struct {
	dgraph.ID
	IsLogin  bool   `json:"isLogin,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateLogin ...
func UpdateLogin(username, oldPassword, newPassword string) bool {
	if CheckLogin(username, oldPassword) {
		login, err := GetHashedPassword(username)
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

// GetHashedPassword ...
func GetHashedPassword(username string) (Login, error) {
	q := `query {
		login(func: has(isLogin)) {
			uid
			username
			password
		}
	}`

	type root struct {
		LoginList []Login `json:"login"`
	}
	newRoot := root{}
	err := executeQuery(q, &newRoot)
	if err != nil {
		return Login{}, err
	}
	return newRoot.LoginList[0], nil
}

// CheckLogin ...
func CheckLogin(username, inputPassword string) bool {
	// get hashed pwd from db
	login, err := GetHashedPassword(username)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return comparePasswords(login.Password, []byte(inputPassword))
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}
