/*
 * licensed to vmware.
 */
package fake

import (
	v1 "apiserver-builder/pkg/client/clientset_generated/clientset/typed/kuber/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKuberV1 struct {
	*testing.Fake
}

func (c *FakeKuberV1) MyKinds(namespace string) v1.MyKindInterface {
	return &FakeMyKinds{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKuberV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
