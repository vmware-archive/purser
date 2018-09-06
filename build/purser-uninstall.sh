# Copyright (c) 2018 VMware Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# change directory to purser directory
cd $HOME/opt/purser-env

# === Purser Controller ===

# Get Cluster config file location
read -p "Taking $HOME/.kube/config as your cluster's configuration (yes/no): " isDefaultConfig
if [ "$isDefaultConfig" != "yes" ]
then
    read -p "Enter location of your cluster's configuration: " kubeConfig
else
    kubeConfig="$HOME/.kube/config"
fi

# uninstall purser controller
echo "Uninstalling purser controller"
kubectl --kubeconfig=$kubeConfig delete -f custom_controller.yaml

# crd uninstall
echo "Uninstalling purser group definition"
kubectl --kubeconfig=$kubeConfig delete -f crd.yaml

rm custom_controller.yaml crd.yaml