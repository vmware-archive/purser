# Manual Installation

To install Purser manually from the Binary follow the steps described below.

## Controller Installation

- Download the purser controller for your environment from the [releases page](https://github.com/vmware/purser/releases/download/v0.1-alpha.2/custom_controller.yaml).

- Install the controller using the `kubectl`.

  ``` bash
  kubectl --kubeconfig=<absolute path to config> create -f custom_controller.yaml
  ```

## Purser Plugin Installation

- Download the purser plugin descriptor for your environment from the [releases page](https://github.com/vmware/purser/releases/download/v0.1-alpha.2/plugin.yaml).

- Move the `plugin.yaml` file into one of the paths specified under the Kubernetes [documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins).

- Download the purser binary corresponding to your operating system from the [releases page](https://github.com/vmware/purser/releases/tag/v0.1-alpha.2).

- Move the binary into one of the directories in your environment `PATH`.
