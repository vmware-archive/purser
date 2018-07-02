/*
 * licensed to vmware.
 */
package v1

import (
	v1 "apiserver-builder/pkg/apis/kuber/v1"
	"apiserver-builder/pkg/client/clientset_generated/clientset/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type KuberV1Interface interface {
	RESTClient() rest.Interface
	MyKindsGetter
}

// KuberV1Client is used to interact with features provided by the kuber.kuber group.
type KuberV1Client struct {
	restClient rest.Interface
}

func (c *KuberV1Client) MyKinds(namespace string) MyKindInterface {
	return newMyKinds(c, namespace)
}

// NewForConfig creates a new KuberV1Client for the given config.
func NewForConfig(c *rest.Config) (*KuberV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &KuberV1Client{client}, nil
}

// NewForConfigOrDie creates a new KuberV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *KuberV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new KuberV1Client for the given RESTClient.
func New(c rest.Interface) *KuberV1Client {
	return &KuberV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *KuberV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
