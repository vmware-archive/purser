# Installation Through Source Code

- [Prerequisites](#prerequisites)
- [Server Side Installation (Controller Installation)](#server_side_installation)
- [Client Side Installation (Plugin Installation)](#client_side_installation)

## Prerequisites

1. Kubernetes Version 1.9 or greater

    - `kubectl` installed and configured. For details refer [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

2. Dependencies

    - [Go](https://golang.org/dl/)

        - version > 1.7
        - setup `GOPATH` environment variable by as per the [Golang documentation](https://github.com/golang/go/wiki/SettingGOPATH).
        - add `$GOPATH/bin` directory to your environment `$PATH` variable.

    - [Docker](https://www.docker.com/get-started)

3. Fetch the Purser source code from GitHub.

   ``` go
   go get github.com/vmware/purser
   ```

   ``` bash
   # change directory to project root
   cd $GOPATH/src/github.com/vmware/purser
   ```

4. For Windows users, install gnu `make` from [here](http://gnuwin32.sourceforge.net/packages/make.htm).

5. Download project dependencies with `make`.

   ``` bash
   # download project tools
   make tools

   # download project dependencies
   make deps

   # update project depedencies
   make update
   ```

## Server Side Installation (Controller Installation)

Follow the below steps to install the purser controller and custom resource definitions for the user groups in the Kubernetes cluster.

### Build Controller Binary

Build the purser controller binary using `make` target.

``` bash
make build
```

### Build Container Image

Update the [Makefile](./Makefile) to set the `REGISTRY` field to your Docker username and execute the following `make` targets to build and publish the docker images.

``` bash
# create the container(docker image)
make container

# authenticate your Docker credentials
docker login

# publish your docker image to docker hub
make push
```

### Install Purser Plugin

- Update the image name in [`purser-controller-setup.yaml`](../cluster/purser-controller-setup.yaml) to the docker image name that you pushed.

- Install the controller in the cluster using `kubectl`.

  ``` bash
  kubectl create -f purser-controller-setup.yaml`
  ```

  _Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

## Client Side Installation (Plugin Installation)

- Build the purser plugin binary in the `GOPATH/bin` directory.

  ``` go
  go build -o $GOPATH/bin/purser_plugin github.com/vmware/purser/cmd/plugin
  ```

- Install the Purser plugin by copying the [`plugin.yaml`](../plugin.yaml) into one of the paths specified under the Kubernetes documentation section [installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).