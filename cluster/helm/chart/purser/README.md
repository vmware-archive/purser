# [Purser](https://github.com/vmware/purser)

Purser is an extension to Kubernetes tasked at providing an insight into cluster topology, costing, capacity allocations and resource interactions along with the provision of logical grouping of resources for Kubernetes based cloud native applications in a cloud neutral manner, with the focus on catering to a multitude of users ranging from Sys Admins, to DevOps to Developers.

It comprises of three components: a controller, a plugin and a UI dashboard.

The controller component deployed inside the cluster watches for K8s native and custom resources associated with the application, thereby, periodically building not just an inventory but also performing discovery by generating and storing the interactions among the resources such as containers, pods and services.

The plugin component is a CLI tool interfacing with the kubectl that helps query costs, savings defined at a level of control of the application level components rather than at the infrastructure level.

The UI dashboard is a robust application that renders the Purser UI for providing visual representation to the complete cluster metrics in a single pane of glass.

> Taken from main [README](https://github.com/vmware/purser/blob/master/README.md)

> [Plugin installation guide](https://github.com/vmware/purser/blob/master/README.md#purser-plugin-setup)

## Chart Configuration

*See `values.yaml` for configuration notes. Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
$ helm install --name purser \
  --set database.storage=10Gi \
    purser
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install --name purser -f values.yaml
```

> **Tip**: You can use the default [values.yaml](values.yaml)