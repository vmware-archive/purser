# Installation through source code

## Prerequisites

1. Kubernetes version 1.9 or greater
    * ``kubectl`` installed and configured. See [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

1. Install dependencies 
    * Install [Go](https://golang.org/dl/)
        - Version atleast 1.7
        - Setup GOPATH environment variable by following [https://github.com/golang/go/wiki/SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)
        - Add $GOPATH/bin directory to your environment $PATH variable
    * Install [Docker](https://www.docker.com/get-started)

1. Get Purser source code
    * `go get github.com/vmware/purser`

1. Change directory to project root
    * `cd $GOPATH/src/github.com/vmware/purser`

1. For windows users, install gnu `make` from [here](http://gnuwin32.sourceforge.net/packages/make.htm)

1. Run the following commands which downloads dependencies
    * `make tools`
    * `make deps`
    * `make update`

## Server side installation

The following two steps installs purser controller and custom resource definitions for user groups in kubernetes cluster.

1. In [Makefile](./Makefile) update `REGISTRY` field to your docker username.

1. Build purser_controller binary using `make build`

1. Create container(docker image) using `make container`

1. Authenticate your docker credentials using `docker login`

1. Push your docker image to docker hub using `make push`

1. In [`custom_controller.yaml`](./custom_controller.yaml) update image name to your docker image name that you pushed

1. Install the controller in the cluster using `kubectl create -f custom_controller.yaml`

## Client side installation

1. Run the following command to create a purser plugin binary in 
   `GOPATH/bin` directory

    `go build -o $GOPATH/bin/purser_plugin github.com/vmware/purser/cmd/plugin`

1. Copy the [plugin.yaml](./plugin.yaml) into one of the paths specified under 
   section [Installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)