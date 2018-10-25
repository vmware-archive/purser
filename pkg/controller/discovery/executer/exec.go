package executer

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

const debug = false

// ExecToPodThroughAPI uninterractively exec to the pod with the command specified.
func ExecToPodThroughAPI(client *kubernetes.Clientset, command, containerName, podName string, stdin io.Reader) (string, string, error) {
	// Prepare the API URL used to execute another process within the Pod. In this case,
	// we'll run a remote shell.
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Command:   strings.Fields(command),
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	if debug {
		log.Debug("Request URL:", req.URL().String())
	}

	config, err := utils.GetKubeconfig("")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch kubeconfig file %v", err)
	}
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	// Connect this process' std{in,out,err} to the remote shell process.
	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}
	return stdout.String(), stderr.String(), nil
}
