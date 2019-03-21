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

package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/vmware/purser/cmd/controller/api/apiHandlers"
	"github.com/vmware/purser/pkg/controller"
)

// StartServer starts api server
func StartServer(cookieStoreKey, cookieName string, conf controller.Config) {
	apiHandlers.SetGroupClient(conf)
	apiHandlers.SetCookieStore(cookieStoreKey, cookieName)
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedCredentials := handlers.AllowCredentials()
	router := NewRouter()
	logrus.Info("Purser server started on port `localhost:3030`")
	logrus.Fatal(http.ListenAndServe(":3030", handlers.CORS(allowedOrigins, allowedCredentials)(router)))
}
