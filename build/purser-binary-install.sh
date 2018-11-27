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
releaseVersion=v1.0.0

# === Purser Plugin ===

# Detecting os type
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac
echo "Detecting your Operating System: ${machine}"

echo "Downloading files for plugin..."
# Download purser plugin yaml
pluginYamlUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/plugin.yaml
wget -q --show-progress -O plugin.yaml $pluginYamlUrl

# Downloading purser plugin binary based on os type
if [ $machine = Linux ]
then
    pluginUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/purser_plugin_linux_amd64
elif [ $machine = Mac ]
then
    pluginUrl=https://github.com/vmware/purser/releases/download/$releaseVersion/purser_plugin_darwin_amd64
else
    echo "No match found for your os: $machine"
    echo "Install the plugin from source code: https://github.com/vmware/purser/blob/master/README.md"
    exit 3  # unsuccessful shell script
fi
wget -q --show-progress -O purser_plugin $pluginUrl

# Move th plugin yaml to one of the location specified in 
# https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/
if [ ! -d $HOME/.kube/plugins ]
then
    mkdir $HOME/.kube/plugins
fi
echo "Moving plugin.yaml to $HOME/.kube/plugins/"
mv plugin.yaml $HOME/.kube/plugins/

# Change execution permissions for the binary
chmod +x purser_plugin

# Move the binary to a location which is in environment PATH variable
echo "Moving the binary to /usr/local/bin"
sudo mv purser_plugin /usr/local/bin

echo "Purser plugin installation Completed"

echo ""

echo "Purser Installation Completed"