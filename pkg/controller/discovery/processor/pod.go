package processor

import (
	"fmt"
	"sync"

	"github.com/vmware/purser/pkg/controller"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/discovery/linker"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var wg sync.WaitGroup

// ProcessPodInteractions fetches details of all the running processes in each container of
// each pod in a given namespace and generates a 1:1 mapping between the communicating pods.
func ProcessPodInteractions(conf *controller.Config) {
	k8sPods := RetrievePodList(conf.Kubeclient, metav1.ListOptions{})

	linker.PopulatePodIPTable(k8sPods)
	processPodDetails(conf.Kubeclient, k8sPods)

	linker.GenerateAndStorePodInteractions()
	log.Infof("Successfully generated Pod To Pod mapping.")
}

func processPodDetails(client *kubernetes.Clientset, pods *corev1.PodList) {
	podsCount := len(pods.Items)
	log.Infof("Processing total of (%d) Pods.", podsCount)

	wg.Add(podsCount)
	{
		for index, pod := range pods.Items {
			log.Debugf("Processing Pod (%d/%d) ... ", index+1, podsCount)

			go func(pod corev1.Pod, index int) {
				defer wg.Done()

				containers := pod.Spec.Containers
				processContainerDetails(client, pod, containers)
				log.Debugf("Finished processing Pod (%d/%d)", index+1, podsCount)

			}(pod, index)
		}
	}
	wg.Wait()
	fmt.Println("")
}
