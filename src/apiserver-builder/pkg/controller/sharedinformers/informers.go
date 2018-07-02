
/*
 * licensed to vmware.
*/


package sharedinformers

// SetupKubernetesTypes registers the config for watching Kubernetes types
func (si *SharedInformers) SetupKubernetesTypes() bool {
    // Set this to true to initial the ClientSet and InformerFactory for
    // Kubernetes APIs (e.g. Deployment)
	return false
}

// StartAdditionalInformers starts watching Deployments
func (si *SharedInformers) StartAdditionalInformers(shutdown <-chan struct{}) {
    // Start specific Kubernetes API informers here.  Note, it is only necessary
    // to start 1 informer for each Kind. (e.g. only 1 Deployment informer)

    // Uncomment this to start listening for Deployment Create / Update / Deletes
    // go si.KubernetesFactory.Apps().V1beta1().Deployments().Informer().Run(shutdown)
}
