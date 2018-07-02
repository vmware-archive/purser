/*
 * licensed to vmware.
 */
package internalversion

import (
	"apiserver-builder/pkg/client/clientset_generated/internalclientset/scheme"
	rest "k8s.io/client-go/rest"
)

type KuberInterface interface {
	RESTClient() rest.Interface
	MyKindsGetter
}

// KuberClient is used to interact with features provided by the kuber.kuber group.
type KuberClient struct {
	restClient rest.Interface
}

func (c *KuberClient) MyKinds(namespace string) MyKindInterface {
	return newMyKinds(c, namespace)
}

// NewForConfig creates a new KuberClient for the given config.
func NewForConfig(c *rest.Config) (*KuberClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &KuberClient{client}, nil
}

// NewForConfigOrDie creates a new KuberClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *KuberClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new KuberClient for the given RESTClient.
func New(c rest.Interface) *KuberClient {
	return &KuberClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	g, err := scheme.Registry.Group("kuber.kuber")
	if err != nil {
		return err
	}

	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != g.GroupVersion.Group {
		gv := g.GroupVersion
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *KuberClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
