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

package apiHandlers

import (
	"encoding/json"
	"encoding/gob"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
)

// Credentials structure
type Credentials struct {
	Password    string `json:"password"`
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

// User structure
type User struct {
	Username string
	Authenticated bool
}

var cookieName = "purser-session-token"

var store *sessions.CookieStore

// initialises cookie store
func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	gob.Register(User{})
}

// LoginUser listens on /auth/login endpoint
func LoginUser(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.Authenticate(cred.Username, cred.Password) {
		logrus.Errorf("wrong credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["user"] = User{
		Username: cred.Username,
		Authenticated: true,
	}

	err = session.Save(r, w)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logrus.Infof("login success")
	w.WriteHeader(http.StatusOK)
}

// LogoutUser listens on /auth/logout endpoint
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// TODO: Enhance
func isUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// ChangePassword listens on /auth/changePassword endpoint
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	addAccessControlHeaders(&w, r)
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.UpdatePassword(cred.Username, cred.Password, cred.NewPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
