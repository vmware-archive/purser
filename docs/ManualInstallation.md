# Manual Installation

* Download the purser controller [here](https://github.com/vmware/purser/releases/download/v0.1-alpha.2/custom_controller.yaml). 
* Install it using `kubectl --kubeconfig=<absolute path to config> create -f custom_controller.yaml`.
* Download purser plugin yaml from [here](https://github.com/vmware/purser/releases/download/v0.1-alpha.2/plugin.yaml).
* Move the plugin.yaml file into one of the paths specified [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins).
* Download the purser binary corresponding to your operating system from [here](https://github.com/vmware/purser/releases/tag/v0.1-alpha.2).
* Move the binary into one of the directory in you environment `PATH`.
