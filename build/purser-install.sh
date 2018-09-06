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
releaseVersion=v0.1-alpha.1

# purser directory
mkdir -p $HOME/opt/purser-env
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

# Download purser controller yaml
echo "Downloading purser controller yaml"
controllerUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/custom_controller.yaml
wget -q --show-progress -O custom_controller.yaml $controllerUrl

# Need crd.yaml if uninstallation is needed
crdUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/crd.yaml
wget -q -O crd.yaml $crdUrl

# Installing purser controller
echo "Installing purser controller"
kubectl --kubeconfig=$kubeConfig create -f custom_controller.yaml

echo ""

# === Purser Plugin ===

# Download purser plugin yaml
echo "Downloading purser plugin yaml"
pluginYamlUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/plugin.yaml
wget -q --show-progress -O plugin.yaml $pluginYamlUrl

# Move th plugin yaml to one of the location specified in 
# https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/
if [ ! -d $HOME/.kube/plugins ]
then
    mkdir $HOME/.kube/plugins
fi
echo "Moving plugin.yaml to $HOME/.kube/plugins/"
mv plugin.yaml $HOME/.kube/plugins/

# Detecting os type
echo "Detecting your os"
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac
echo ${machine}

# Downloading purser plugin binary based on os type
echo "Downloading purser plugin binary"
if [ $machine = Linux ]
then
    # echo "Downloading from https://github.com/vmware/purser/releases/download/v0.1-alpha.1/purser_plugin_linux_amd64"
    pluginUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/purser_plugin_linux_amd64
elif [ $machine = Mac ]
then
    # echo "Downloading from https://github.com/vmware/purser/releases/download/v0.1-alpha.1/purser_plugin_darwin_amd64"
    pluginUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/purser_plugin_darwin_amd64
else
    echo "No match found for your os: $machine"
    echo "Install the plugin from source code: https://github.com/vmware/purser/blob/master/README.md"
    exit 3  # unsuccessful shell script
fi
wget -q --show-progress -O purser_plugin $pluginUrl

# Change execution permissions for the binary
chmod +x purser_plugin

# Move the binary to a location which is in environment PATH variable
echo "Moving the binary to /usr/local/bin"
sudo mv purser_plugin /usr/local/bin

echo "Installation Completed"