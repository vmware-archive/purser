# Usage

Once installed, Purser is ready for use right away. You can query using native Kubernetes grouping artifacts.

Following are the commands that purser supports.
Use flag `--kubeconfig=<absolute path to config>` if your cluster configuration is not at [default location].(https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable).

```
kubectl plugin purser get summary
kubectl plugin purser get savings
kubectl plugin purser get resources namespace <Namespace>
kubectl plugin purser get resources label <key=val>
kubectl plugin purser get cost label <key=val>
kubectl plugin purser get cost pod <pod name>
kubectl plugin purser get cost node <node name>
kubectl plugin purser set user-costs
kubectl plugin purser get user-costs
```

**Examples:**


1. Get cluster summary


        $ kubectl plugin purser get summary
            Cluster Summary
            Compute:
                Node count:                 57
                Cost:                       3015.48$
                Total Capacity:
                    Cpu(vCPU):               456
                    Memory(GB):              1770.50
                Provisioned Resources:
                    Cpu Request(vCPU):       319
                    Memory Request(GB):      1032.67
            Storage:
                Persistent Volume count:    151
                Capacity(GB):               9297.00
                Cost:                       4124.79$
                PV Claim count:             108
                PV Claim Capacity(GB):      8867.00
            Cost:
                Compute cost:               3015.48$
                Storage cost:               4124.79$
                Total cost:                 7140.27$



2. Get cost of all nodes

        kubectl purser get cost node all

3. Get savings

        $ kubectl plugin purser get savings
            Savings Summary
            Storage:
                Unused Volumes:             43
                Unused Capacity(GB):        430.00
                Month To Date Savings:      186.33$
                Projected Monthly Savings:   1066.40$


Next, define higher level groupings to define your business, logical or application constructs

### Defining custom groups
Group .yaml format

```
kind: Group
metadata:
    name: <name of the group>
spec:
    labels:
        <label1>
        ....
        <labelN>
    namespace:
        <namespace1,...namespaceN>
```
**Example:**

Query the cost of Cost Insight infrastructure deployed in "default" namespace

1. The following is the ci.yaml definition which groups a few native Kubernetes labels into a business/application construct

```
apiVersion: vmware.kuber/v1
kind: Group
metadata:
  name: ci-group
spec:
  name: ci-group
  labels:
    app: vrbc-transformer
    app: vrbc-adapterdefault
    ```
2. Create the construct defined above

        kubectl create -f ci.yaml

3. Get the cost of CI group

        kubectl get cost group CI
