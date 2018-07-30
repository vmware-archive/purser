package client

import (
	"kuber-controller/crd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type GroupCrdClient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

// This file implement all the (CRUD) client methods we need to access Group CRD object
func CreateGroupCrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string) *GroupCrdClient {
	return &GroupCrdClient{cl: cl, ns: namespace, plural: crd.GroupPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

func (f *GroupCrdClient) CreateGroup(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

/*func (f *GroupCrdClient) CreateSubscriber(obj *crd.Subscriber) (*crd.Subscriber, error) {
	var result crd.Subscriber
	err := f.cl.Post().
		Namespace(f.ns).Resource(crd.SubscriberPlural).
		Body(obj).Do().Into(&result)
	return &result, err
}*/

func (f *GroupCrdClient) UpdateGroup(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.
	//Put().
		Put().Name((obj.Name)).
	//Patch(types.JSONPatchType).Name(obj.Name).
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *GroupCrdClient) DeleteGroup(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *GroupCrdClient) GetGroup(name string) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *GroupCrdClient) ListGroups(opts meta_v1.ListOptions) (*crd.GroupList, error) {
	var result crd.GroupList
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

/*func (f *SubscriberCrdClient) ListSubscribers(opts meta_v1.ListOptions) (*crd.SubscriberList, error) {
	var result crd.SubscriberList
	err := f.cl.Get().
		Namespace(f.ns).Resource(crd.SubscriberPlural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}*/

// Create a new List watch for our TPR
func (f *GroupCrdClient) NewListWatchGroup() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}
