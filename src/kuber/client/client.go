package client

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"kuber/crd"
)

// This file implement all the (CRUD) client methods we need to access our CRD object

func CrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string) *Crdclient {
	return &Crdclient{cl: cl, ns: namespace, plural: crd.CRDPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

type Crdclient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

func (f *Crdclient) Create(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Update(obj *crd.Group) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.
		//Put().
		Put().Name((obj.Name)).
		//Patch(types.JSONPatchType).Name(obj.Name).
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) Get(name string) (*crd.Group, error) {
	var result crd.Group
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) List(opts meta_v1.ListOptions) (*crd.GroupList, error) {
	var result crd.GroupList
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}
