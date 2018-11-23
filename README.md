# Purser extension for K8s

[![Build Status](https://travis-ci.org/vmware/purser.svg?branch=master)](https://travis-ci.org/vmware/purser) [![Go Report Card](https://goreportcard.com/badge/github.com/vmware/purser)](https://goreportcard.com/report/github.com/vmware/purser)

- [What is Purser?](#purser)
- [Features](#features)
- [Setup and Installation](#setup-and-installation)
- [Uninstalling](#uninstalling)
- [Additional Documentation](#additional-documentation)
- [Community, Discussion, Contribution and Support](#community-discussion-contribution-and-support)

## Purser

Purser is an extension to Kubernetes tasked at providing an insight into *cluster topology*, *costing*, and *resource interactions* along with the provision of *logical grouping of resources* for Kubernetes based cloud native applications in a cloud neutral manner, with the focus on catering to a multitude of users ranging from Sys Admins, to DevOps to Developers.

It comprises of three components: a controller, a plugin and a UI dashboard.  

The controller component deployed inside the cluster watches for K8s native and custom resources associated with the application, thereby, periodically building not just an inventory but also performing discovery by generating and storing the interactions among the resources such as containers, pods and services.

The plugin component is a CLI tool interfacing with the `kubectl` that helps query costs, savings defined at a level of control of the application level components such as the _Memory and CPU consumptions and utilization_ rather than at the infrastructure level.

The UI dashboard is a robust application that renders the Purser UI for providing visual representation to the complete cluster metrics in a single pane of glass. 

### Demo

![demo](/docs/img/purser-cli.gif)

## Features

Purser with it's robust CLI and UI capabilities provides a set of features including, but not limited to the list below.
 
 - Capability to provide visibility into the following aspects of the K8s cluster
    - workload cost associated with the native/custom resources
    - savings opportunities associated with storage and compute requirements
    - single pane view of the complete cluster hierarchy
    - interactions among associated resources such as pods and services
 
 - Capability of user defined logical grouping of resources based on `K8s CRD` implementation for enhanced filtering.
 
 - A plugin extension to `kubectl` along with the UI for developer centric usage.

### Purser UI demo

 ![demo](/docs/img/purser-ui.gif)

## Setup and Installation

Follow the instructions below to set up Purser in your environment.  

### Prerequisites

- Kubernetes version 1.9 or greater.
- `kubectl` installed and configured. For details see [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

### Installation

The preferred and the quickest way to install purser is through Binary installation.

#### OS-specific installation methods

#### Linux and macOS

``` bash
# Binary installation
wget -q https://github.com/vmware/purser/releases/download/v0.1-alpha.2/purser-install.sh && sh purser-install.sh
```

Enter your cluster's configuration path when prompted. The plugin binary needs to be in your `PATH` environment variable, so once the download of the binary is finished the script tries to move it to `/usr/local/bin`. This may need your sudo permission.

#### Windows

For installation on Windows follow the steps in the [manual installation guide](./docs/manual-installation.md).

#### Other Installation Methods

For other installation methods such as **manual installation** or **installation from source code** refer guides in [docs](./docs).

## Uninstalling

### Linux/macOS

``` bash
wget -q https://github.com/vmware/purser/releases/download/v0.1-alpha.2/purser-uninstall.sh && sh purser-uninstall.sh
```

### Others

``` bash
kubectl delete -f custom_controller.yaml
kubectl delete -f crd.yaml
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

## Additional Documentation

Additional documentation can be found below:

- [Manual Installation Guide](https://github.com/vmware/purser/blob/master/docs/manual-installation.md)
- [Source Code Installation Guide](https://github.com/vmware/purser/blob/master/docs/sourcecode-installation.md)
- [Purser Architecture and Workflow](https://github.com/vmware/purser/blob/master/docs/architecture.md)
- [Purser Plugin Usage](https://github.com/vmware/purser/blob/master/docs/plugin-usage.md)
- [Developers Guide](https://github.com/vmware/purser/blob/master/docs/developers-guide.md)
- [Purser Deployment Guide](https://github.com/vmware/purser/blob/master/docs/purser-deployment.md)
- [Purser UI Development Guide](https://github.com/vmware/purser/blob/master/ui/README.md)

## Community, Discussion, Contribution and Support

**Issues:** Have an issue with Purser, please [log it](https://github.com/vmware/purser/issues).

**Contributing:** Would you like to contribute to our project, refer [How to contribute](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md) docs.
