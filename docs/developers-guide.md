# Developers Guide

- [Prerequisites](#prerequisites)
- [Workspace Setup](#workspace-setup)
- [Database Setup](#database-setup)
- [Running Purser Controller](#running-purser-controller)
- [Running Purser UI](#running-purser-ui)
- [Purser Plugin Compilation](#purser-plugin-compilation)
- [Plugin Execution](#plugin-execution)

## Prerequisites

1. Ensure the following dependencies are installed on your system.

   - [Go](https://golang.org/dl/)
   - [Git](https://git-scm.com/downloads)
   - [Docker](https://www.docker.com/)

   You may use the official binaries or your usual package manager.
   Also set the following environment variables
   - Set `GOPATH` environment variable. Refer [setting GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)
   - Add `$GOPATH/bin` in system `PATH` variable by running `export PATH=$PATH:$GOPATH/bin`.

      Optionally, add the above exports to your `.bash_profile` or `.bashrc` to persist across console sessions.

2. Verify that the dependencies are properly installed.

   ``` bash
   go version, should be at least 1.7

   git version

   docker version
   ```

## Workspace Setup

### Fork the repository
Navigate to the [Purser repo on GitHub](https://github.com/vmware/purser) and use the 'Fork' button. 
This gives you a copy of the repo for pull requests back to purser in `https://github.com/<your-github-id>/purser`

### Clone and Set Upstream Remote

Make a local clone of the forked repo and add the base purser
repo as the upstream remote repository.

``` shell
# create and change directory to $GOPATH/src/github.com/vmware
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware

# clone the forked repository and change directory to purser
git clone https://github.com/<your-github-id>/purser.git
cd purser

# add upstream repository as the original purser repo
git remote add upstream https://github.com/vmware/purser.git
```

The last git command prepares your clone to pull changes from the
upstream repo and push them into the fork, which enables you to keep
the fork up to date.

### Download dependencies

Run the following commands to download dependencies.

``` shell
make tools
make deps
make install
```

## Database Setup

In order to persist inventory and discovery information such as pods and service details we use
[Dgraph](https://dgraph.io/) to store the inventory metrics and resource relationship.

In order to install DGraph from docker image follow the following steps:

- Pull the latest Dgraph version

  ```bash
  docker pull dgraph/dgraph
  ```

- To run Dgraph in Docker

  ```bash
  mkdir -p /tmp/data
  
  # Run dgraph-zero
  docker run -d -p 5080:5080 -p 6080:6080 -p 8080:8080 -p 9080:9080 -p 8000:8000 -v /tmp/data:/dgraph --name diggy dgraph/dgraph dgraph zero

  # In another terminal, now run dgraph-alpha
  docker exec -d diggy dgraph alpha --lru_mb 2048 --zero localhost:5080
  ```
 
 - Optional: To start Dgraph UI(at `localhost:8000`) for running manual queries
   ```bash
   # Run Dgraph Ratel
   docker exec -d diggy dgraph-ratel
   ```
   
## Running Purser Controller
To run purser controller execute following commands

```bash
# change directory to purser main folder
cd $GOPATH/src/github.com/vmware/purser

# run purser with log level as info and interactions as disabled by default
go run cmd/controller/purserctrl.go --kubeconfig=<path-to-your-cluster-config> --interactions=disable --dgraphURL=localhost --log=info
```

## Running Purser UI
Install latest version of `node` and `npm`. Then to run purser UI execute the following commands
```bash
# change directory to purser ui folder
cd $GOPATH/src/github.com/vmware/purser/ui

# install node modules
npm install

# run purser UI at localhost:4200
npm run startdev
```
_Refer [UI docs](../ui/README.md) for more details._

## Purser Plugin Compilation

To create purser plugin binary `purser_plugin` at path `$GOPATH/bin` run the following commands
  ```bash
  # change directory to purser main folder
  cd $GOPATH/src/github.com/vmware/purser
  
  # create binary at path $GOPATH/bin
  go build -o $GOPATH/bin/purser_plugin github.com/vmware/purser/cmd/plugin
  ```

**NOTE:** _Windows users need to rename `purser_plugin` to `purser_plugin.exe`_

## Plugin Execution

1. In order to install the Purser plugin, copy the [plugin.yaml](../plugin.yaml) file to one of the specified paths defined under the section [installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

2. Run the following command to check the purser plugin works locally.

   ``` bash
   kubectl --kubeconfig=<absolute path to kubeconfig file> plugin purser help
   ```

## Useful commands and links
- To contribute to purser refer [CONTRIBUTING](../CONTRIBUTING.md) and [CODE_OF_CONDUCT](../CODE_OF_CONDUCT.md)
- To drop complete dgraph database: `curl -X POST localhost:8080/alter -d '{"drop_all": true}'`