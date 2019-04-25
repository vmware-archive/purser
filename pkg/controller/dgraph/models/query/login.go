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

// Authenticate performs user authentication for service access
func Authenticate(username, inputPassword string) bool {
	if !validateUsername(username) {
		return false
	}
	login, err := getLoginCredentials(username)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return comparePasswords(login.Password, []byte(inputPassword))
}

// UpdatePassword updates stored password with new one for the given username in Dgraph
func UpdatePassword(username, oldPassword, newPassword string) bool {
	if Authenticate(username, oldPassword) {
		login, err := getLoginCredentials(username)
		if err != nil {
			logrus.Error(err)
			return false
		}
		if err = hashAndUpdatePassword(&login, newPassword); err == nil {
			return true
		}
		logrus.Error(err)
	}
	return false
}

func hashAndUpdatePassword(login *dgraph.Login, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}
	login.Password = string(hashedPassword)
	_, err = dgraph.MutateNode(login, dgraph.UPDATE)
	return err
}

// getLoginCredentials returns a struct of hashed password and username.
func getLoginCredentials(username string) (dgraph.Login, error) {
	q := `query {
		login(func: has(isLogin)) @filter(eq(username, ` + username + `)) {
			uid
			username
			password
		}
	}`
	type root struct {
		LoginList []dgraph.Login `json:"login"`
	}
	newRoot := root{}
	if err := executeQuery(q, &newRoot); err != nil || newRoot.LoginList == nil {
		return dgraph.Login{}, err
	}
	return newRoot.LoginList[0], nil
}

func validateUsername(username string) bool {
	return username == "admin"
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	if err := bcrypt.CompareHashAndPassword(byteHash, plainPwd); err != nil {
		logrus.Error(err)
		return false
	}
	return true
}
