# Purser extension for K8s
[![Build Status](https://travis-ci.org/vmware/purser.svg?branch=master)](https://travis-ci.org/vmware/purser)

Cost visbility for Kubernetes based Cloud Native Applications

## Why?

Today, cost visibility in the world of Cloud Native Applications is very limited. It is mostly restricted to cost of cloud infrastructure at a
high level and usually involves a lot of manual steps or custom scripting.

Wouldn't it be great if you know the cost of you Kubernetes deployed applications, not matter the cloud of your choice? Don't you wish there was an easy way to
incorporate your budgeting and cost savings at a level of control that was entirely based on application level components rather than infrastructure? 

## What is Purser

Purser provides cost visibility of services, microservices and applications deployed with Kubernetes in a cloud neutral manner. It does so at a granular level and
over time ranges that match with budget planning.

Purser is an extension to Kubernetes. More specifically, it is a tool interfacing with ``kubectl`` that helps you query for cost based on native Kubernetes artifacts
as well as your own custom defined services. In addition, Purser allows for alerting on budget adherence and helps enforce budgets and savings.

Purser currently supports Kubernetes deployments on Amazon Web Services. Support for VMware vSphere, Azure, Google Compute Engine are planned.


## Features

* Query cost associated with Kubernetes native groups
* Extend Purser with YAML based declarative custom service, microservice and application definitions
* Capability for control over time range for cost query
* Capability for cost analysis based on resource Usage or Allocation
* Visibility into Cost savings oppurtunities
* Set budget limits on Kubernetes native or custom defined groups
* Capability to enforce budget for Kubernetes native or custom defined groups

## Getting Started

Instructions to install and start using Purser plugin.

### Installation

#### Prerequisites

* Kubernetes version 1.8 or greater
* ``kubectl`` installed and configured. See [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

#### Installation through binaries

*Note: Installation through binaries is in progress. Follow next section Installation through source code*
##### Server side installation

The following two steps installs purser controller and custom resource definitions for user groups in kubernetes cluster.

1. Installing purser custom controller
    * wget https://github.com/vmware/purser/blob/master/custom_controller.yaml
    * kubectl apply -f custom_controller.yaml

2. Installing CRDs(Custom Resource Definitions) for custom groups.
    * wget https://github.com/vmware/purser/blob/master/crd.yaml
    * kubectl create -f crd.yaml
    
    Note: The above CRD is also created by purser custom controller, if CRD is already by controller then kubectl displays resource already exist message.

##### Client side installation

The following two steps installs the necessary components on client side.

1. Downloading kubectl plugin yaml file
    * wget https://github.com/vmware/purser/blob/master/plugin.yaml
    * copy the plugin.yaml file into one of the paths specified in `Plugin loader` section in [link](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

2. Installing kubectl plugin binary
    * Follow [CODECOMPILE.md](./docs/CODECOMPILE.md)

#### Installation through source code

1. Install dependencies 
    * Install [Go](https://golang.org/dl/)
        - Version atleast 1.7
        - Setup GOPATH environment variable by following [https://github.com/golang/go/wiki/SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)
    * Install [Docker](https://www.docker.com/get-started)

1. Get Purser source code
    * `go get github.com/vmware/purser`

1. Change directory to project root
    * `cd $GOPATH/src/github.com/vmware/purser`

##### Server side installation

The following two steps installs purser controller and custom resource definitions for user groups in kubernetes cluster.

1. In [Makefile](./Makefile) update `REGISTRY` field to your docker username.

1. Build purser_controller binary using `make build`

1. Create container(docker image) using `make container`

1. Authenticate your docker credentials using `docker login`

1. Push your docker image to docker hub using `make push`

1. In kubernetes cluster download custom_controller.yaml from [here](https://github.com/vmware/purser/blob/master/custom_controller.yaml) or

    `wget https://github.com/vmware/purser/blob/master/custom_controller.yaml`

1. In `custom_controller.yaml` update image name to your docker image name that you pushed

1. Install the controller in the cluster using `kubectl create -f custom_controller.yaml`

##### Client side installation

1. Run the following command to create a purser plugin binary in 
   `GOPATH/bin` directory

    `go install github.com/vmware/purser/cmd/purser_plugin`

1. Copy the [plugin.yaml](../plugin.yaml) into one of the paths specified under 
   section [Installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)


### Usage

Once installed, Purser is ready for use right away. You can query using native Kubernetes grouping artifacts

**Examples:**


1. Get cost of pods having label "app=vrbc-adapter"


        $ kubectl purser get cost label app=vrbc-adapter
            ===Pods Cost Details===
            Pod Name:                     vrbc-adapter-statefulset-1-1-577-0
            Node:                         ip-172-20-40-248.ec2.internal
            Pod Compute Cost Percentage:  7.03
            Persistent Volume Claims:     
                vrbc-adapter-volume-1-1-577-vrbc-adapter-statefulset-1-1-577-0
            Cost:                         
            Total Cost:          108.092667$
            Compute Cost:        69.426000$
            Storage Cost:        38.666667$

            Pod Name:                     vrbc-adapter-statefulset-1-1-577-1
            Node:                         ip-172-20-41-91.ec2.internal
            Pod Compute Cost Percentage:  6.96
            Persistent Volume Claims:     
                vrbc-adapter-volume-1-1-577-vrbc-adapter-statefulset-1-1-577-1
            Cost:                         
                Total Cost:          107.412371$
                Compute Cost:        68.745704$
                Storage Cost:        38.666667$

            Pod Name:                     vrbc-adapter-statefulset-1-1-577-2
            Node:                         ip-172-20-52-245.ec2.internal
            Pod Compute Cost Percentage:  5.86
            Persistent Volume Claims:     
                vrbc-adapter-volume-1-1-577-vrbc-adapter-statefulset-1-1-577-2
            Cost:                         
                Total Cost:          96.496567$
                Compute Cost:        57.829900$
                Storage Cost:        38.666667$
                
            Total Cost Summary:           
                Total Cost:          312.001604$
                Compute Cost:        196.001604$
                Storage Cost:        116.000000$


2. Get cost of all nodes

        kubectl purser get cost node all


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






