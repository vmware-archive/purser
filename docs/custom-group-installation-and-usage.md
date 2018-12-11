# Custom Group Installation and Usage

To get resource and cost visibility for a particular set of pods Purser allows user to create custom logical group.
User can define the label filter logic while creating the logical group i.e, pods satisfying these conditions will belong to this custom group.

## Installing logical group definition and an example logical group

To install the logical group definition into your cluster, 
download [purser-group-crd.yaml](../cluster/artifacts/purser-group-crd.yaml) yaml i.e,
```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: groups.vmware.purser.com
spec:
  group: vmware.purser.com
  names:
    kind: Group
    listKind: GroupList
    plural: groups
    singular: group
  scope: Namespaced
  version: v1
status:
  acceptedNames:
    kind: Group
    listKind: GroupList
    plural: groups
    singular: group
```
and use kubectl to install this definition
```bash
kubectl create -f purser-group-crd.yaml
```
_**NOTE:** This installation is needed only once per cluster_

**Installing an example logical group**

Download [example-group.yaml](../cluster/artifacts/example-group.yaml) yaml i.e,
```yaml
apiVersion: vmware.purser.com/v1
kind: Group
metadata:
  name: example-group
spec:
  name: example-group
  labels:
    expr1:
      app:
        - sample-app
        - sample-app2
      env:
        - dev
    expr2:
      namespace:
        - ns1
        - ns2
    expr3:
      key1:
        - val1
      key2:
        - val2
```
and use kubectl to create this logical group
```bash
kubectl create -f example-group.yaml
kubectl get groups.vmware.purser.com
```

This will create a custom logical group with name `example-group` of type `groups.vmware.purser.com`.
The label filter (used to fetch pods belonging to this group) for `example-group` will be
```yaml
(app=sampl-app OR app=sample-app2 OR env=dev) AND (namespace=ns1 OR namespace=ns2) AND (key1=val1 OR key2=val2)
```

In general the syntax purser supports is:

```
expr1 AND expr2 AND expr3 AND ...
where each expr is of form key1:value1 OR key2:value2 OR key1:value3 OR ...
```

## Usage
For resource and cost visibility into this newly created logical group run the following command
```bash
kubectl plugin purser get resources group example-group
```
_Refer [purser installation](../README.md#installation) to install purser controller and plugin_ 

## Uninstalling purser custom group
To uninstall purser custom group run the following command
```bash
kubectl delete -f purser-group-crd.yaml
```
where [purser-group-crd.yaml](../cluster/artifacts/purser-group-crd.yaml) is same file that you downloaded during installation.