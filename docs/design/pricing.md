# Pricing in Purser

## User defined pricing
Currently purser supports user defined pricing for cpu, memory and storage resources per hour. The default pricing for the resources are

* CPU: 0.024$  per vCPU per Hour
* Memory:  0.01$ per GB per Hour
* Storage:  0.00013888888$ per GB per Hour

These default pricing for cpu and memory have been set by taking average [prices](https://aws.amazon.com/ec2/pricing/on-demand/) in AWS ec2 instances.
For storage we set pricing proportional to 0.1$ per GB per month referring to AWS [ebs pricing](https://aws.amazon.com/ebs/pricing/).

User can edit these pricing using purser plugin: `kubectl plugin purser set user-costs`

_(Future work)_ Option to edit these default pricing in UI.

## Using node labels for accurate pricing (WIP)
Data needed to get correct pricing of node:

* Cloud Provider (AWS, GCE etc)
* Region (us-east etc)
* Machine Type (t2.micro, m4.large)
* Operating system (linux, windows etc)
* Rate card which gives cost of node depending on above data
* _(Future work)_ Costing based on instance type i.e., "Is the instance On-demand or Spot-Instance or Reserved-Instance etc?"
* _(Future work)_ Discounts

### Getting cloud provider, region and machine type
Kubelet populates few [reserved labels](https://kubernetes.io/docs/reference/kubernetes-api/labels-annotations-taints/#beta-kubernetes-io-instance-type) on nodes. Using these labels we can determine region, machineType and operating system.
Command `kubectl describe node <nodeName>` gives labels.

* Default assume cloud provider as aws. Check section [Finding Cloud Provider](#finding-cloud-provider).
* Label `beta.kubernetes.io/instance-type=m4.10xlarg` gives machine type. Here for this example machineType is m4.10xlarge
* Label `beta.kubernetes.io/os=linux` gives operating system. Here os is linux
* Label `failure-domain.beta.kubernetes.io/region=us-west-1` gives region. Here region is us-west-1

_Note: kubelet will not set these reserved labels if the cluster is not using cloud provider._

If any of the required labels is not available then we should fall back to default pricing.



#### Storage Volume Pricing

`kubectl get pv` gives us storage class(ex: gp2, my-storage-class etc) for each volume.

Further using command `kubectl describe storageclass <storageclass-name>` will give output in which there will be a field `Parameters`
 containing labels for `type` of storage (Ex: gp2) and `zone` (Ex: us-west-1c).

_Parameters_ field may not contain _zone_ label. In such case we can get region from `kubectl describe pv <pv-name>` using label `failure-domain.beta.kubernetes.io/region`.

If _type_ label is also not present then we should fall back to default pricing.



### Finding Cloud Provider
While initiating a cluster either by kubeadm or kops or other kubernetes installers the user will set cloud-provider, if it isn't set kubernetes assumes that cluster is being deployed on bare metal. 
Further when a new node is created `.spec.providerID` will be set based (by _kubelet_) on cloud-provider.
The value of `providerId` will be `(providerName +  "://" + instanceID)`. As we need `providerName` (aws, azure etc) we can get it from `providerID`.
If getting cloud provider name fails we should fallback and assume default prices for aws.

*Example cluster on aws: --cloud-provider=aws command-line flag is needed (to successfully register the node with cloud provider) to be present for the API server, controller manager, and every kubelet in the cluster.

References:

* kubeadm: https://kubernetes.io/docs/concepts/cluster-administration/cloud-providers/
* providerID value: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/cloud-provider/cloud.go#L94
* Cluster on aws: https://blog.heptio.com/setting-up-the-kubernetes-aws-cloud-provider-6f0349b512bd
* More on providerID: https://blog.scottlowe.org/2018/09/28/setting-up-the-kubernetes-aws-cloud-provider/


*Note: All above kubectl commands will have corresponding methods in kubernetes client-go



### Populating Rate Card
#### Design

* Identify cloud provider, region of the cluster
* Embed the crawler code in purser and run cloud specific crawler to fetch the rate cards.
* Populate the data in dgraph.
* Update rate card periodically.
* Support: AWS, Azure, PKS, VKE, GCE

#### AWS:

* Public API: Available
* Reference: https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/price-changes.html
* API call: https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/region/index.json
* Example for us-east-1: https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/us-east-1/index.json
* Note: aws provides sdk in golang for pricing. Reference: https://docs.aws.amazon.com/sdk-for-go/api/service/pricing/