package processor

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RetrievePodList returns list of pods in the given namespace.
func RetrievePodList(client *kubernetes.Clientset, options metav1.ListOptions) *corev1.PodList {
	pods, err := client.CoreV1().Pods("").List(options)
	if err != nil {
		log.Errorf("failed to retrieve pods: %v", err)
	}
	return pods
}
