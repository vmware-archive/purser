package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	getPodsCommand = "kubectl --kubeconfig=%s get pods -l app=%s -o json"
	getNodeCommand = "kubectl --kubeconfig=%s get node %s -o json"
)

func getNodeDetails(nodeName string) Node {
	node := Node{}
	command := fmt.Sprintf(getNodeCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), nodeName)
	bytes := executeCommand(command)
	json := string(bytes)
	node.name = nodeName
	node.instanceType = gjson.Get(json, "metadata.labels.beta\\.kubernetes\\.io/instance-type").Str
	return node
}

func getPodsForLabel(label string) []Pod {
	pods := []Pod{}
	command := fmt.Sprintf(getPodsCommand, os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_KUBECONFIG"), label)
	bytes := executeCommand(command)
	json := string(bytes)
	items := gjson.Get(json, "items")

	items.ForEach(func(key, value gjson.Result) bool {
		name := value.Get("metadata.name")
		nodeName := value.Get("spec.nodeName")
		pod := Pod{name: name.Str, nodeName: nodeName.Str}
		pods = append(pods, pod)
		return true
	})
	return pods
}

func executeCommand(command string) []byte {
	slice := strings.Fields(command)
	cmd := exec.Command(slice[0], slice[1:]...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}

func printPodDetails(pods []Pod) {
	fmt.Println("===POD Details===")
	fmt.Println("POD Name \t\t\t\t\t Node Name")
	for _, value := range pods {
		fmt.Println(value.name + " \t" + value.nodeName)
	}
}

func printNodeDetails(nodes []Node) {
	fmt.Println("===Node Details===")
	fmt.Println("Node Name \t\t\t InstanceType")
	for _, value := range nodes {
		fmt.Println(value.name + " \t" + value.instanceType)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func testNodeDescribe() {
	dat, err := ioutil.ReadFile("/Users/gurusreekanthc/sources/kuber-plugin/node.output")
	check(err)
	parseNodeDescribe(dat)
	//fmt.Print(string(dat))
}

func main() {
	/*inputs := os.Args[1:]
	pods := getPodsForLabel(inputs[1])
	printPodDetails(pods)
	node := getNodeDetails(pods[0].nodeName)
	printNodeDetails([]Node{node})*/
	testNodeDescribe()
}
