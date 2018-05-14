# Kuber plugin for K8s

Cost visbility for Kubernetes based Container Native Applications

## Why?

Today, cost visibility in the world of Container Native Applications is very limited. It is mostly restricted to cost of cloud infrastructure at a
high level and usually involves a lot of manual steps or custom scripting.

Wouldn't it be great if you know the cost of you Kuberentes deployed applications, not matter the cloud of your choice? Don't you wish there was an easy way to
incorporate your budgeting and cost savings at a level of control that was entirely based on application level components rather than infrastructure? 

## What is Kuber

Kuber provides cost visibility of services, microservices and applications deployed with Kubernetes in a cloud neutral manner. It does so at a granular level and
over time ranges that match with budget planning.

Kuber is a CLI extension to Kubernetes. More specifically, it is a kubectl plugin that helps you query for cost based on native Kubernetes artifacts
as well as your own custom defined services. In addition, kuber allows for alerting on budget adherence and helps enforce budgets and savings.

## How does it work

Using Kuber is simple and similar to the declarative philosphy adpated by Kubeernetes and kubectl
```
kubectl plugin kuber get_cost label app=my-web-ui
```
Would get you the monthly aggregated cost of kubernetes PODs labeled with 'app=my-web-ui'
```
kubectl plugin kuber set_limit namespace bakend-auto-scaling-group month limit 2400 action alert email backend devops@org.com
```
Would set a monthly limit of 2400$ on resources in the 'backend-auto-scaling-group' and if the cost this group breaches the limit, send an email to devops@org.com
with cost details

## Features

* Query cost associated with Kubernetes native groups
* Extend Kuber with YAML based declarative custom service, microservice and applicatin definitions
* Control over time range for cost query
* Set budget limits on Kubernetes native or custom defined groups
* Cost analysis based on reosurce Usage or Allocation

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites


### Installing




