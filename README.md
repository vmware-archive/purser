# Kuber extension for K8s

Cost visbility for Kubernetes based Cloud Native Applications

## Why?

Today, cost visibility in the world of Cloud Native Applications is very limited. It is mostly restricted to cost of cloud infrastructure at a
high level and usually involves a lot of manual steps or custom scripting.

Wouldn't it be great if you know the cost of you Kuberentes deployed applications, not matter the cloud of your choice? Don't you wish there was an easy way to
incorporate your budgeting and cost savings at a level of control that was entirely based on application level components rather than infrastructure? 

## What is Kuber

Kuber provides cost visibility of services, microservices and applications deployed with Kubernetes in a cloud neutral manner. It does so at a granular level and
over time ranges that match with budget planning.

Kuber is an extension to Kubernetes. More specifically, it is a tool interfacing with ``kubectl`` that helps you query for cost based on native Kubernetes artifacts
as well as your own custom defined services. In addition, Kuber allows for alerting on budget adherence and helps enforce budgets and savings.

Kuber currently supports Kubernetes deployments on Amazon Web Services. Support for VMware vSphere, Azure, Google Compute Engine are planned.


## Features

* Query cost associated with Kubernetes native groups
* Extend Kuber with YAML based declarative custom service, microservice and application definitions
* Capability for control over time range for cost query
* Capability for cost analysis based on resource Usage or Allocation
* Visibility into Cost savings oppurtunities
* Set budget limits on Kubernetes native or custom defined groups
* Capability to enforce budget for Kubernetes native or custom defined groups

## Getting Started

Instructions to install and start using Kuber plugin.

### Prerequisites

* Kubernetes version 1.8 or greater
* ``kubectl`` installed and configured. See [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Installing

Installing is as simple as downloading and executing a shell script: [kuber_install.sh] ()

#### What does the installation do?

Installing Kuber creates a few Kubernetes supported extensions on your cluster. This enables install once, query from anywhere using ``kubectl``

1. Installing CRDs(Custom Resource Definitions)
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/crd.yaml
    * kubectl create -f crd.yaml
2. Installing custom controller
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/custom_controller.yaml
    * kubectl apply -f custom_controller.yaml
3. Installing API extension server
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/api_extension_server.yaml
    * kubectl apply -f api_extension_server.yaml


### Usage

Once installed, Kuber is ready for use right away. You can query using native Kubernetes grouping artifacts

**Examples:**


1. Get cost of pods having label "app=heimdall"

        kubectl kuber get cost label app=heimdall

2. Get cost of all nodes

        kubectl kuber get cost node all


Next, define higher level groupings to define your business, logical or application constructs

### Defining custom groups
Group .yaml format

```
Kind: Group
Metadata:
    name: <name of the group>
Spec:
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
    Kind: Group
    Metadata:
        name: CI
    Spec:
        labels:
            app=vrbc-transformer
            app=vrbc-adapter
            app=vrbc-showback
            app=vrbc-ui
            app=ci-lambda
        namespace:
            default
    ```
2. Create the construct defined above

        kubectl create -f ci.yaml

3. Get the cost of CI group

        kubectl get cost group CI

### Uninstalling

Not convinced? Uninstalling cleans up everything Kuber and leaves your cluster in it's original state: [kuber_uninstall.sh] ()

## Enabling historic cost

## Utilization based cost






