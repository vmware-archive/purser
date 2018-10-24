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

package client

import (
	"flag"
	"os/user"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/utils"

	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

const environment = "prod"

// GetAPIExtensionClient returns a client for the cluster and it's config.
func GetAPIExtensionClient() (*apiextcs.Clientset, *rest.Config) {
	var config *rest.Config
	var err error

	if environment == "dev" {
		kubeconfig := flag.String("kubeconf", getUsrHomeDir()+"/.kube/config", "path to the kubeconfig file")
		flag.Parse()
		config, err = utils.GetKubeconfig(*kubeconfig)
		if err != nil {
			log.Fatalf("failed to fetch kubeconfig %v", err)
		}
	} else {
		config, err = utils.GetKubeconfig("")
		if err != nil {
			log.Fatalf("failed to fetch kubeconfig %v", err)
		}
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to connect to the cluster %v", err)
	}

	return clientset, config
}

func getUsrHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to fetch path to config file %v", err)
	}
	return usr.HomeDir
}
