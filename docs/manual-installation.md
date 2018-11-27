# Manual Installation

To install Purser manually from the Binary follow the steps described below.

## Purser Controller Setup
Download the controller setup yaml file from [here](https://github.com/vmware/purser/blob/master/cluster/purser-controller-setup.yaml).
``` bash
# Controller installation
kubectl create -f purser-controller-setup.yaml
```

## Purser UI Setup
Download the UI setup yaml file from [here](https://github.com/vmware/purser/blob/master/cluster/purser-ui-setup.yaml).
``` bash
# UI installation
kubectl create -f purser-ui-setup.yaml
```

## Purser Plugin Installation

- Download the purser plugin descriptor for your environment from the [releases page](https://github.com/vmware/purser/releases/download/v1.0.0/plugin.yaml).

- Move the `plugin.yaml` file into one of the paths specified under the Kubernetes [documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins).

- Download the purser binary corresponding to your operating system from the [releases page](https://github.com/vmware/purser/releases/tag/v1.0.0).

- Move the binary into one of the directories in your environment `PATH`.
