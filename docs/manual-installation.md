# Manual Installation

To install Purser manually from the Binary follow the steps described below.

## Purser Setup
The following steps will install Purser in your cluster at namespace `purser`.
Creation of this namespace is needed because purser needs to create a service-account which requires namespace.
Also, the frontend will use kubernetes DNS to call backend for data and this DNS contains a field for namespace.
``` bash
# Namespace setup
kubectl create ns purser

# DB setup
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/purser-database-setup.yaml -O
kubectl --namespace=purser create -f purser-database-setup.yaml

# Purser controller setup
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/purser-controller-setup.yaml -O
kubectl --namespace=purser create -f purser-controller-setup.yaml

# Purser UI setup
curl https://raw.githubusercontent.com/vmware/purser/master/cluster/purser-ui-setup.yaml -O
kubectl --namespace=purser create -f purser-ui-setup.yaml
```
**NOTE:** If you don't have `curl` installed you can download `purser-database-setup.yaml` from [here](./cluster/purser-database-setup.yaml), `purser-controller-setup.yaml` from [here](cluster/purser-controller-setup.yaml) and `purser-ui-setup.yaml` from [here](cluster/purser-ui-setup.yaml). 
Then `kubectl create -f purser-database-setup.yaml` ,
`kubectl create -f purser-controller-setup.yaml` and `kubectl create -f purser-ui-setup.yaml` will setup purser in your cluster.

##### Change Settings and Enable/Disable Purser Features

The following settings can be customized before Controller installation:

- Change the default **log level**, **dgraph url** and **dgraph port** by editing `args` field in the [purser-controller-setup.yaml](cluster/purser-controller-setup.yaml). (Default: `--log=info`, `--dgraphURL=purser-db`, `--dgraphPort=9080`)
- Enable/Disable **resource interactions** capability by editing `args` field in the [purser-controller-setup.yaml](cluster/purser-controller-setup.yaml) and uncommenting `pods/exec` rule from purser-permissions. (Default: `disabled`)
- Enable **subscription to inventory changes** capability by creating an object of custom resource kind `Subscriber`. (Refer: [example-subscriber.yaml](./cluster/artifacts/example-subscriber.yaml))
- Enable **customized logical grouping of resources** by creating an object of custom resource kind `Group`. (Refer: [docs](docs/custom-group-installation-and-usage.md) for custom group installation and usage)

_**NOTE:** Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at the [default location](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable)._

## Purser Plugin Installation

- Download the purser plugin descriptor for your environment from the [releases page](https://github.com/vmware/purser/releases/download/v1.0.0/plugin.yaml).

- Move the `plugin.yaml` file into one of the paths specified under the Kubernetes [documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins).

- Download the purser binary corresponding to your operating system from the [releases page](https://github.com/vmware/purser/releases/tag/v1.0.0).

- Move the binary into one of the directories in your environment `PATH`.
