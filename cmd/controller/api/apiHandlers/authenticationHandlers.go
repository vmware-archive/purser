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
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/vmware/purser/pkg/controller/dgraph/models/query"
)

// Credentials structure
type Credentials struct {
	Password    string `json:"password"`
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

var cookieName = "purser-session-token"

var store sessions.Store

// SetCookieStore initialises cookie store
func SetCookieStore(cookieKey, cookiename string) {
	store = sessions.NewCookieStore([]byte(cookieKey))
	cookieName = cookiename
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

	if !query.CheckLogin(cred.Username, cred.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	saveSession(session, w, r)
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
	session.Values["authenticated"] = false
	saveSession(session, w, r)
}

func saveSession(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	err := session.Save(r, w)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func isUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, cookieName)
	if err != nil {
		logrus.Errorf("unable to get session from cookie store, err: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return false
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}
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

	if !query.UpdateLogin(cred.Username, cred.Password, cred.NewPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
