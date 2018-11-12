# Architecture of Purser

The following diagram represents the architecture of Purser.

![Architecture](/docs/img/architecture.png)

The following are the main componenets installed in Kubernetes for Purser.

1. **Kubernetes API Server Extension**

    All the Purser `kubectl` commands hit the API server extension. These APIs understand the input command, compute and return the required output.

2. **Custom Controller**

    The custom controller watches for changes in state of pods, nodes, persistent volumes, etc. and update the inventory in CRDs.

3. **Custom Resource Definitions(CRDs)**

    Custom Resource Definitions are like any other resource(Pod, Node, etc.) and store the config data like `Group Definitions` and inventory.

4. **Metric Store**

    Metric store is used to store the utilization, allocation metrics of inventory and also calculated costs.

5. **CRON Job**

    CRON Job collects the stats of inventory and calculates the cost periodically and stores in Metric Store.

## Work Flow

1. Purser installation steps create Custom Controller, CRON Job and CRDs in Kubernetes.

2. Once installed the custom controller collects all the inventory(pods, nodes, pv, etc.) and stores in CRDs, later it watches for any changes in inventory and stores the changes in CRDs.

3. CRON Job kicks in periodically and collect the stats and stores the stats in metric store. CRON Job also calculates the Costs in the same cycle and stores them in the metric store.

4. Any `kubectl` command invocations are received by Kubernetes API server extension.  APIs then process the required output based on the configurations(for groups), inventory, costs metrics and returns to the user.
