# Purser extension for K8s

[![Build Status](https://travis-ci.org/vmware/purser.svg?branch=master)](https://travis-ci.org/vmware/purser) [![Go Report Card](https://goreportcard.com/badge/github.com/vmware/purser)](https://goreportcard.com/report/github.com/vmware/purser)

- [What is Purser?](#purser)
- [Features](#features)
- [Setup and Installation](#setup-and-installation)
- [Uninstalling](#uninstalling)
- [API Documentation](#api-documentation)
- [Additional Documentation](#additional-documentation)
- [Community, Discussion, Contribution and Support](#community-discussion-contribution-and-support)

## Purser

Purser is an extension to Kubernetes tasked at providing an insight into *cluster topology*, *costing*, *capacity allocations* and *resource interactions* along with the provision of *logical grouping of resources* for Kubernetes based cloud native applications in a cloud neutral manner, with the focus on catering to a multitude of users ranging from Sys Admins, to DevOps to Developers.

It comprises of three components: a controller, a plugin and a UI dashboard.  

The controller component deployed inside the cluster watches for K8s native and custom resources associated with the application, thereby, periodically building not just an inventory but also performing discovery by generating and storing the interactions among the resources such as containers, pods and services.

The plugin component is a CLI tool interfacing with the `kubectl` that helps query costs, savings defined at a level of control of the application level components  rather than at the infrastructure level.

The UI dashboard is a robust application that renders the Purser UI for providing visual representation to the complete cluster metrics in a single pane of glass. 

### Demo

![demo](/docs/img/purser-cli.gif)

## Features

Purser with it's robust CLI and UI capabilities provides a set of features including, but not limited to the list below.
 
 - Capability to provide visibility into the following aspects of the K8s cluster
    - workload cost associated with the native/custom resources
    - savings opportunities associated with storage and compute requirements
    - single pane view of the complete cluster hierarchy
    - capacity allocations for CPU, memory, disk space and other resources
    - interactions among associated resources such as pods and services
 
 - Capability of user defined logical grouping of resources based on `K8s CRD` implementation for enhanced filtering.
 
 - A plugin extension to `kubectl` along with the UI for developer centric usage.
 
 - Capability to notify inventory changes via web-hook implementation. 

### Purser UI demo

 ![demo](https://user-images.githubusercontent.com/42761785/49430222-74d25600-f7d0-11e8-97ad-ba1388fb6d8f.gif)

## Setup and Installation

Follow the instructions below to set up Purser in your environment.  

### Prerequisites

- Kubernetes version 1.9 or greater.
- `kubectl` installed and configured. For details see [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

### Installation

Purser has three components to install.

- [Purser Controller Setup](./README.md#Purser-Controller-Setup)
- [Purser UI Setup](./README.md#Purser-UI-Setup)
- [Purser Plugin Setup](./README.md#Purser-Plugin-Setup)

#### Purser Controller Setup
Download the controller setup yaml file from [here](./cluster/purser-controller-setup.yaml).

**To enable/disable Purser features**

Edit [purser-controller-setup.yaml](./cluster/purser-controller-setup.yaml)
* Choose **log level** by editing `args` of purser-controller deployment (default: info)
* Enable/Disable discovery of **interactions** feature by editing `args` of purser-controller deployment 
and uncommenting `pods/exec` rule from purser-permissions (default: disabled)
* Change **dgraph's** url and port number by editing `args` of purser-controller deployment (default: purser-db, 9080)
``` bash
# Controller installation
kubectl create -f purser-controller-setup.yaml
```

#### Purser UI Setup
Download the UI setup yaml file from [here](./cluster/purser-ui-setup.yaml).
``` bash
# UI installation
kubectl create -f purser-ui-setup.yaml
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

#### Purser Plugin Setup

#### Linux and macOS

``` bash
# Binary installation
wget -q https://github.com/vmware/purser/blob/master/build/purser-binary-install.sh && sh purser-binary-install.sh
```

Enter your cluster's configuration path when prompted. The plugin binary needs to be in your `PATH` environment variable, so once the download of the binary is finished the script tries to move it to `/usr/local/bin`. This may need your sudo permission.

#### Windows

For installation on Windows follow the steps in the [manual installation guide](./docs/manual-installation.md).

#### Other Installation Methods

For other installation methods such as **manual installation** or **installation from source code** refer guides in [docs](./docs).

## Uninstalling

### Uninstalling Purser Controller
``` bash
kubectl delete -f purser-controller-setup.yaml
kubectl delete pvc datadir-dgraph-0
```

### Uninstalling Purser UI
``` bash
kubectl delete -f purser-ui-setup.yaml
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

### Uninstalling Purser Binary

### Linux/macOS

``` bash
wget -q https://github.com/vmware/purser/blob/master/build/purser-binary-uninstall.sh && sh purser-binary-uninstall.sh
```

## API Documentation

The project uses Swagger to document API's endpoints. The documentation is available at [Swagger Hub](https://app.swaggerhub.com/apis/hemani19/purser/1.0.0).

## Additional Documentation

Additional documentation can be found below:

- [Manual Installation Guide](./docs/manual-installation.md)
- [Source Code Installation Guide](./docs/sourcecode-installation.md)
- [Purser Architecture and Workflow](./docs/architecture.md)
- [Purser Plugin Usage](./docs/plugin-usage.md)
- [Developers Guide](./docs/developers-guide.md)
- [Purser Deployment Guide](./docs/purser-deployment.md)
- [Purser UI Development Guide](./ui/README.md)

## Community, Discussion, Contribution and Support

**Issues:** Have an issue with Purser, please [log it](https://github.com/vmware/purser/issues).

**Contributing:** Would you like to contribute to our project, refer [How to contribute](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md) docs.
