# Copyright (c) 2018 VMware Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#    http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Realease Version
releaseVersion=1.0.2

echo "Installing Purser (minimal setup) version: ${releaseVersion}"

# Namespace setup
echo "Creating namespace purser"
kubectl create ns purser

# DB setup
echo "Setting up database for Purser"
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/minimal/purser-database-setup.yaml -O
kubectl --namespace=purser create -f purser-database-setup.yaml
echo "Waiting for database containers to be in running state... (30s)"
sleep 30s

# Purser controller setup
echo "Setting up controller for Purser"
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/minimal/purser-controller-setup.yaml -O
kubectl --namespace=purser create -f purser-controller-setup.yaml

# Purser UI setup
echo "Setting up UI for Purser"
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/minimal/purser-ui-setup.yaml -O
kubectl --namespace=purser create -f purser-ui-setup.yaml

echo "Purser setup is completed"
