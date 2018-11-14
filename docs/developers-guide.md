# Developers Guide

- [Prerequisites](prerequisites)
- [Database Setup](database_setup)
- [Binary Compilation](binary_compilation)
- [Local Execution](local_execution)

## Prerequisites

1. Ensure the following dependencies are installed on your system.

   - [Go](https://golang.org/dl/)
   - [Git](https://git-scm.com/downloads)
   - [Docker](https://www.docker.com/)

   You may use the official binaries or your usual package manager.

2. Verify that the dependencies are properly installed.

   ``` bash
   go version, should be at least 1.7

   git version

   docker version
   ```

## Database Setup

In order to persist inventory and discovery information such as pods and service details we use
[Dgraph](https://dgraph.io/) to store the relationship and other metrics.

In order to install DGraph from docker image follow the following steps:

- Pull the latest Dgraph version

  ```bash
  docker pull dgraph/dgraph
  ```

- To run Dgraph in Docker

  ```bash
  mkdir -p ~/dgraph
  

  # Run dgraphzero
  docker run -it -p 5080:5080 -p 6080:6080 -p 8080:8080 -p 9080:9080 -p 8000:8000 -v ~/dgraph:/dgraph --name dgraph dgraph/dgraph dgraph zero

  # In another terminal, now run dgraph
  docker exec -it dgraph dgraph server --lru_mb 2048 --zero localhost:5080

  # And in another, run ratel (Dgraph UI)
  docker exec -it dgraph dgraph-ratel
  ```

## Binary Compilation

1. If the `GOPATH` environment variable isn't set then set it as per the documentation [here](https://github.com/golang/go/wiki/SettingGOPATH).

2. Add `GOPATH/bin` to your `PATH` environment variable.

   ``` bash
   export PATH=$PATH:$GOPATH/bin
   ```

   Optionally, add the above exports to your `.bash_profile` to persist across console sessions.

3. Fetch and install the Purser project.

   ``` go
   go get github.com/vmware/purser
   ```

   Run the following command to create a purser plugin binary in `GOPATH/bin` directory.

   ``` go
   go install github.com/vmware/purser/cmd/purser_plugin
   ```

## Local Execution

1. In order to install the Purser plugin, copy the [plugin.yaml](../plugin.yaml) file to one of the specified paths defined under the section [installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

2. Run the following command to check the purser plugin works locally.

   ``` bash
   kubectl --kubeconfig=<absolute path to kubeconfig file> plugin purser help
   ```