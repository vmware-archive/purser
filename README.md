# Kuber plugin for K8s

Cost visbility for Kubernetes based Cloud Native Applications

## Why?

Today, cost visibility in the world of Cloud Native Applications is very limited. It is mostly restricted to cost of cloud infrastructure at a
high level and usually involves a lot of manual steps or custom scripting.

Wouldn't it be great if you know the cost of you Kuberentes deployed applications, not matter the cloud of your choice? Don't you wish there was an easy way to
incorporate your budgeting and cost savings at a level of control that was entirely based on application level components rather than infrastructure? 

## What is Kuber

Kuber provides cost visibility of services, microservices and applications deployed with Kubernetes in a cloud neutral manner. It does so at a granular level and
over time ranges that match with budget planning.

Kuber is a CLI extension to Kubernetes. More specifically, it is a ``kubectl`` plugin that helps you query for cost based on native Kubernetes artifacts
as well as your own custom defined services. In addition, kuber allows for alerting on budget adherence and helps enforce budgets and savings.

Kuber currently supports Kubernetes deployments on Amazon Web Services. Support for VMware vSphere, Azure, Google Compute AEngine and other platforms are planned.

## How does it work

Using Kuber is simple and similar to the declarative philosphy adpated by Kubeernetes and kubectl
```
kuber get_cost label app=my-web-ui
```
Would get you the monthly aggregated cost of kubernetes PODs labeled with 'app=my-web-ui'
```
kuber set_limit namespace backend-auto-scaling-group month limit 2400 action alert email backend devops@org.com
```
Would set a monthly limit of 2400$ on resources in the 'backend-auto-scaling-group' and if the cost this group breaches the limit, send an email to devops@org.com
with cost details

## Features

* Query cost associated with Kubernetes native groups
* Extend Kuber with YAML based declarative custom service, microservice and application definitions
* Capability for control over time range for cost query
* Capacbility for cost analysis based on resource Usage or Allocation
* Set budget limits on Kubernetes native or custom defined groups
* Capability to enforce budget and cost saving for Kubernetes native or custom defined groups

## Getting Started

Instructions to install and start using Kuber plugin.

### Prerequisites

You must have ``kubectl`` installed and configured. See [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Installing

1. Installing CRDs(Custom Resource Definitions)
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/crd.yaml
    * kubectl create -f crd.yaml
2. Installing custom controller
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/custom_controller.yaml
    * kubectl create -f custom_controller.yaml
3. Installing API extension server
    * wget https://gitlab.eng.vmware.com/kuber/kuber-plugin/blob/master/api_extension_server.yaml
    * kubectl create -f api_extension_server.yaml


### Usage


**kubectl --kubeconfig=<path to kubeconfig> kuber get cost {label|pod|node|group} <_variable> [duration=hourly|weekly|monthly]**

**Examples:**


1. Get cost of pods having label "app=heimdall"

        kubectl --kubeconfig=/Users/abc/prod kuber get cost label app=heimdall

2. Get cost of all nodes

        kubectl --kubeconfig=/Users/abc/prod kuber get cost node all

3. Get weekly cost of a pod

        kubectl --kubeconfig=/Users/abc/prod kuber get cost pod pod-name

## Advanced Usage

Users can create advanced groups using kubernetes CRDs(Custom Resource Definitions) and query the cost using group name.

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

1. CRD file definition
    
    The following is the ci.yaml definition.

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
2. Create the CRD in kubernetes

        kubectl --kubeconfig=/Users/abc/prod create -f ci.yaml

3. Get the cost of CI group.

        kubectl --kubeconfig=/Users/abc/prod kuber get cost group CI






