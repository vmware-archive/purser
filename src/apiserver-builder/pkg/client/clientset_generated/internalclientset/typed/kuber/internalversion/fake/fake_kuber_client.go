/*
 * licensed to vmware.
 */
package fake

import (
	internalversion "apiserver-builder/pkg/client/clientset_generated/internalclientset/typed/kuber/internalversion"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKuber struct {
	*testing.Fake
}

func (c *FakeKuber) MyKinds(namespace string) internalversion.MyKindInterface {
	return &FakeMyKinds{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKuber) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
