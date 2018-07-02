
/*
 * licensed to vmware.
*/


package mykind

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"apiserver-builder/pkg/apis/kuber/v1"
	"apiserver-builder/pkg/controller/sharedinformers"
	listers "apiserver-builder/pkg/client/listers_generated/kuber/v1"
)

// +controller:group=kuber,version=v1,kind=MyKind,resource=mykinds
type MyKindControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about MyKind
	lister listers.MyKindLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *MyKindControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing mykinds labels
	c.lister = arguments.GetSharedInformers().Factory.Kuber().V1().MyKinds().Lister()
}

// Reconcile handles enqueued messages
func (c *MyKindControllerImpl) Reconcile(u *v1.MyKind) error {
	// Implement controller logic here
	log.Printf("Running reconcile MyKind for %s\n", u.Name)
	return nil
}

func (c *MyKindControllerImpl) Get(namespace, name string) (*v1.MyKind, error) {
	return c.lister.MyKinds(namespace).Get(name)
}
