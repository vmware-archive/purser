![logo](https://user-images.githubusercontent.com/42761785/53145168-2f4e4980-35c5-11e9-867b-8d637671ec23.png)
# K8s Extension for Application Visibility

[![Build Status](https://travis-ci.org/vmware/purser.svg?branch=master)](https://travis-ci.org/vmware/purser) [![Go Report Card](https://goreportcard.com/badge/github.com/vmware/purser)](https://goreportcard.com/report/github.com/vmware/purser)

- [What is Purser?](#overview)
- [Features](#features)
- [Setup and Installation](#setup-and-installation)
- [Uninstalling](#uninstalling)
- [API Documentation](#api-documentation)
- [Additional Documentation](#additional-documentation)
- [Community, Discussion, Contribution and Support](#community-discussion-contribution-and-support)

## Overview

Purser is an extension to Kubernetes tasked at providing an insight into *cluster topology*, *costing*, *capacity allocations* and *resource interactions* along with the provision of *logical grouping of resources* for Kubernetes based cloud native applications in a cloud neutral manner, with the focus on catering to a multitude of users ranging from Sys Admins, to DevOps to Developers.

It comprises of three components: a controller, a plugin and a UI dashboard.  

The controller component deployed inside the cluster watches for K8s native and custom resources associated with the application, thereby, periodically building not just an inventory but also performing discovery by generating and storing the interactions among the resources such as containers, pods and services.

The plugin component is a CLI tool interfacing with the `kubectl` that helps query costs, savings defined at a level of control of the application level components  rather than at the infrastructure level.

The UI dashboard is a robust application that renders the Purser UI for providing visual representation to the complete cluster metrics in a single pane of glass. 

## Features

Purser with its robust CLI and UI capabilities provides a set of features including, but not limited to the list below.
 
 - Capability to provide visibility into the following aspects of the K8s cluster
    - workload cost associated with the native/custom resources
    - savings opportunities associated with storage and compute requirements
    - single pane view of the complete cluster hierarchy
    - capacity allocations for CPU, memory, disk space and other resources
    - interactions among associated resources such as pods and services
 
 - Capability of user defined logical grouping of resources based on `K8s CRD` implementation for enhanced filtering.
 
 - A plugin extension to `kubectl` along with the UI for developer centric usage.
 
 - Capability to subscribe to inventory changes via web-hook implementation. 

### UI Demo

 ![demo](https://user-images.githubusercontent.com/42461220/54865566-35367680-4d8d-11e9-9e07-e9aa77d7c6ec.gif)

### CLI Demo

 ![demo](/docs/img/purser-cli.gif)

## Setup and Installation

### Prerequisites
- Kubernetes version 1.9 or greater.
- `kubectl` installed and configured. For details see [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).
- `curl` installed. Download it [here](https://curl.haxx.se/download.html)

### Linux/Mac Users:
```bash
curl https://raw.githubusercontent.com/vmware/purser/master/build/purser-setup.sh -O && sh purser-setup.sh
```

### Windows/Other Users:

For detailed installation steps follow the instructions in the [manual installation guide](./docs/manual-installation.md).

### Default Login
To login default username is `admin` and password is `purser!123`.
_NOTE: Please change password after first login.

### Purser Plugin Setup (Optional)
_NOTE: This Plugin installation is optional._

If you want to install and use Purser's command line interface
- [Plugin installation guide](./docs/plugin-installation.md).
- [Plugin Usage](./docs/plugin-usage.md).


### Other Installation Methods

For other installation methods such as **manual installation** or **installation from source code** refer guides in [docs](./docs).

### Uninstalling

``` bash
kubectl delete ns purser
```

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._


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

**Contributing:** Would you like to contribute to our project, refer [How to contribute](./CONTRIBUTING.md), [Developers Guide](./docs/developers-guide.md) and [Code of Conduct](./CODE_OF_CONDUCT.md) docs.
