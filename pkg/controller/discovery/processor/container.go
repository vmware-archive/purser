package processor

import (
	"fmt"
	"strings"

	"github.com/vmware/purser/pkg/controller/discovery/linker"

	log "github.com/Sirupsen/logrus"
	"github.com/vmware/purser/pkg/controller/discovery/executer"
	"github.com/vmware/purser/pkg/controller/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const debug = false

func processContainerDetails(client *kubernetes.Clientset, pod corev1.Pod, containers []corev1.Container) {
	for _, container := range containers {
		pidList, cmdList := getPIDList(client, container.Name, pod.Name)
		for index, pid := range pidList {
			process := linker.Process{ID: pid, Name: cmdList[index]}
			getProcessDump(client, pod, container.Name, process)
		}
	}
}

func getPIDList(client *kubernetes.Clientset, containerName, podName string) ([]string, []string) {
	command := "ps -A -o pid,cmd"
	output, err := executeCommandInPod(client, command, containerName, podName)
	if err != nil {
		return nil, nil
	}

	pidCMDList := strings.Split(output, "\n")
	var pidList, cmdList []string

	for _, pidCMD := range pidCMDList {
		if pidCMD != "" {
			pidCMDClean := strings.Split((strings.TrimSpace(pidCMD)), " ")
			pidList = append(pidList, pidCMDClean[0])
			cmdList = append(cmdList, strings.Join(pidCMDClean[1:], " "))
		}
	}
	// ignore first line i.e, PID CMD headers
	return pidList[1:], cmdList[1:]
}

func getProcessDump(client *kubernetes.Clientset, pod corev1.Pod, containerName string, process linker.Process) {
	//get tcp information from /proc/pid/net/tcp for each process
	if process.ID != "" {
		tcpCommand := "cat /proc/" + process.ID + "/net/tcp"
		tcpOutput, err := executeCommandInPod(client, tcpCommand, containerName, pod.Name)
		if err != nil {
			//to clean dump only to have required fields
			tcpDump := utils.PurgeTCPData(tcpOutput)
			linker.PopulateMappingTables(tcpDump, pod, containerName)
		}

		tcp6Command := "cat /proc/" + process.ID + "/net/tcp6"
		tcp6Output, err := executeCommandInPod(client, tcp6Command, containerName, pod.Name)
		if err != nil {
			//to clean dump only to have required fields
			tcp6Dump := utils.PurgeTCP6Data(tcp6Output)
			linker.PopulateMappingTables(tcp6Dump, pod, containerName)
		}
	}
}

func executeCommandInPod(client *kubernetes.Clientset, command, containerName, podName string) (string, error) {
	output, stderr, err := executer.ExecToPodThroughAPI(client, command, containerName, podName, nil)

	if err != nil && debug {
		log.Debugf("Failed `exec`ing to the container %q, command %q Error: %+v", podName, command, err)
	}

	if len(stderr) > 0 {
		log.Warnf("stderr: %v", stderr)
		err = fmt.Errorf("stderr: %v", stderr)
	}

	return output, err
}
