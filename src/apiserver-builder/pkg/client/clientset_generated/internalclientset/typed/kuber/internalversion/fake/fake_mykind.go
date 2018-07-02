/*
 * licensed to vmware.
 */
package fake

import (
	kuber "apiserver-builder/pkg/apis/kuber"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMyKinds implements MyKindInterface
type FakeMyKinds struct {
	Fake *FakeKuber
	ns   string
}

var mykindsResource = schema.GroupVersionResource{Group: "kuber.kuber", Version: "", Resource: "mykinds"}

var mykindsKind = schema.GroupVersionKind{Group: "kuber.kuber", Version: "", Kind: "MyKind"}

// Get takes name of the myKind, and returns the corresponding myKind object, and an error if there is any.
func (c *FakeMyKinds) Get(name string, options v1.GetOptions) (result *kuber.MyKind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mykindsResource, c.ns, name), &kuber.MyKind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kuber.MyKind), err
}

// List takes label and field selectors, and returns the list of MyKinds that match those selectors.
func (c *FakeMyKinds) List(opts v1.ListOptions) (result *kuber.MyKindList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mykindsResource, mykindsKind, c.ns, opts), &kuber.MyKindList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kuber.MyKindList{}
	for _, item := range obj.(*kuber.MyKindList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested myKinds.
func (c *FakeMyKinds) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mykindsResource, c.ns, opts))

}

// Create takes the representation of a myKind and creates it.  Returns the server's representation of the myKind, and an error, if there is any.
func (c *FakeMyKinds) Create(myKind *kuber.MyKind) (result *kuber.MyKind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mykindsResource, c.ns, myKind), &kuber.MyKind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kuber.MyKind), err
}

// Update takes the representation of a myKind and updates it. Returns the server's representation of the myKind, and an error, if there is any.
func (c *FakeMyKinds) Update(myKind *kuber.MyKind) (result *kuber.MyKind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mykindsResource, c.ns, myKind), &kuber.MyKind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kuber.MyKind), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMyKinds) UpdateStatus(myKind *kuber.MyKind) (*kuber.MyKind, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mykindsResource, "status", c.ns, myKind), &kuber.MyKind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kuber.MyKind), err
}

// Delete takes name of the myKind and deletes it. Returns an error if one occurs.
func (c *FakeMyKinds) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mykindsResource, c.ns, name), &kuber.MyKind{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMyKinds) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mykindsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &kuber.MyKindList{})
	return err
}

// Patch applies the patch and returns the patched myKind.
func (c *FakeMyKinds) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kuber.MyKind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mykindsResource, c.ns, name, data, subresources...), &kuber.MyKind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kuber.MyKind), err
}
