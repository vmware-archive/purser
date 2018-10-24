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

package main

import (
	"github.com/vmware/purser/cmd/controller/config"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/discovery/processor"
	"github.com/vmware/purser/pkg/controller/eventprocessor"
	"github.com/vmware/purser/pkg/utils"
)

var conf controller.Config

func init() {
	utils.InitializeLogger()
	config.Setup(&conf)
}

func main() {
	go eventprocessor.ProcessEvents(&conf)
	processor.ProcessPodInteractions(&conf)
	controller.Start(&conf)
}
