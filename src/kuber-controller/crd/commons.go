package crd

import (
	"reflect"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCRD(clientset apiextcs.Interface, fullName string, group string, version string, plural string) error {
	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: fullName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   group,
			Version: version,
			//TODO: make cluster scoped?
			Scope: apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Plural: plural,
				Kind:   reflect.TypeOf(Group{}).Name(),
			},
		},
	}
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	// Ignore error if it already exists
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
