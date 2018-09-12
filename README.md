# Purser extension for K8s
[![Build Status](https://travis-ci.org/vmware/purser.svg?branch=master)](https://travis-ci.org/vmware/purser)

Cost visbility for Kubernetes based Cloud Native Applications.

## Why?

Today, the cost visibility in the world of Cloud Native Applications is very limited. It is mostly restricted to cost of cloud 
infrastructure at a high level and usually involves a lot of manual steps or custom scripting.

Wouldn't it be great if you know the cost of you Kubernetes deployed applications, not matter the cloud of your choice? Don't 
you wish there was an easy way to incorporate your budgeting and cost savings at a level of control that was entirely based on 
application level components rather than infrastructure? 

## What is Purser

Purser provides cost visibility of services, microservices and applications deployed with Kubernetes in a cloud neutral 
manner. It does so at a granular level and over time ranges that match with budget planning.

Purser is an extension to Kubernetes. More specifically, it is a tool interfacing with ``kubectl`` that helps you query for 
cost based on native Kubernetes artifacts as well as your own custom defined services. In addition, Purser allows for alerting 
on budget adherence and helps enforce budgets and savings.

Purser currently supports Kubernetes deployments on Amazon Web Services. Support for VMware vSphere, Azure, Google Compute 
Engine are planned.

## Features

* Query cost associated with Kubernetes native groups.
* Extend Purser with YAML based declarative custom service, microservice and application definitions.
* Capability for control over time range for cost query.
* Capability for cost analysis based on resource Usage or Allocation.
* Visibility into Cost savings opportunities.
* Set budget limits on Kubernetes native or custom defined groups.
* Capability to enforce budget for Kubernetes native or custom defined groups.

## Use Case

Currenty the below list of commands are supported for the Purser plugin. 

``` bash
# Query cluster visibility in terms of savings and summary for the application. 
kubectl plugin purser get [summary|savings]

# Query resources filtered by associated namespace, labels and groups.
kubectl plugin purser get resources namespace <Namespace>
kubectl plugin purser get resources label <key=val>
kubectl plugin purser get resources group <group-name>

# Query cost filtered by associated labels, pods and node.
kubectl plugin purser get cost label <key=val>
kubectl plugin purser get cost pod <pod name>
kubectl plugin purser get cost node <node name>

# Configure user-costs for the choice of deployment.
kubectl plugin purser [set|get] user-costs
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>`, if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)_.

For detailed usage with examples see [here](./docs/Usage.md).

## Installation

### Prerequisites

* Kubernetes version 1.9 or greater.
* ``kubectl`` installed and configured. See [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

### Installation Methods

* [Binary (Preferred method)](#Binary-Installation)
* [Manual Installation](./docs/ManualInstallation.md)
* [Source Code](./docs/SourceCodeInstallation.md)

### Binary Installation

#### Linux and macOS:

``` bash
wget -q https://github.com/vmware/purser/releases/download/v0.1-alpha.2/purser-install.sh && sh purser-install.sh
```

Enter your cluster's configuration path when prompted. We need the plugin binary to be in your `PATH` environment variable, so 
once the download of the binary is finished the script tries to move it to `/usr/local/bin`. This may need your sudo 
permission.

#### Windows:

Windows users, follow the steps under [manual installation](./docs/ManualInstallation.md) section.

### Manual Installation

Refer [manual installation docs](./docs/ManualInstallation.md).

### Source Code

For detailed installation throught source code, refer [this](./docs/SourceCodeInstallation.md).

## Uninstallation

**For Linux and Mac Users:**

``` bash
wget -q https://github.com/vmware/purser/releases/download/v0.1-alpha.2/purser-uninstall.sh && sh purser-uninstall.sh
```

**For Others:**

``` bash
kubectl delete -f custom_controller.yaml
kubectl delete -f crd.yaml
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

## Contributors

For developers who would like to contribute to our project refer [How to contribute](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md) docs.
