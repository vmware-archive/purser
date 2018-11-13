package graph

import (
	"github.com/vmware/purser/pkg/controller/dgraph/models"
	"strconv"
)

// Node represents each node in the graph
// ID: unique id of pod
// Label: pod name
// Title: string "pods"
// Value: number of times pod has communicated with others
// Group: Connected component number, used for coloring different components in different colors
// CID: list of all services the pod belongs to.
type Node struct {
	ID       int      `json:"id"`
	Label    string   `json:"label"`
	Title    string   `json:"title"`
	Value    int      `json:"value"`
	Group    int      `json:"group"`
	Contains []string `json:"contains"`
}

// Edge represents each edge in the graph
// From: unique id of source pod
// TO: unique id of destination pod
// Title: string containing number of times these two pods communicated
type Edge struct {
	From  int    `json:"from"`
	To    int    `json:"to"`
	Title string `json:"title"`
}

// GetPodNodesAndEdges ...
func GetPodNodesAndEdges(pods []models.Pod) ([]Node, []Edge) {
	uniqueIDs, numConnections := getPodUniqueIDsAndNumConnections(pods)
	nodes := createPodNodes(pods, uniqueIDs, numConnections)
	edges := createPodEdges(pods, uniqueIDs)
	return nodes, edges
}

func getPodUniqueIDsAndNumConnections(pods []models.Pod) (map[string]int, map[string]int) {
	uniqueIDs := make(map[string]int)
	numConnections := make(map[string]int)
	for _, pod := range pods {
		setPodUniqueIDsAndNumConnections(pod, uniqueIDs, numConnections)
	}
	return uniqueIDs, numConnections
}

func setPodUniqueIDsAndNumConnections(pod models.Pod, uniqueIDs, numConnections map[string]int) {
	uniqueID := len(uniqueIDs) + 1
	if _, isPresent := uniqueIDs[pod.Name]; !isPresent {
		uniqueIDs[pod.Name] = uniqueID
		numConnections[pod.Name] = 0
		for _, dstPod := range pod.Interacts {
			numConnections[pod.Name] += int(dstPod.Count)
		}
	}
}

func createPodNodes(pods []models.Pod, uniqueIDs, numConnections map[string]int) []Node {
	nodes := []Node{}
	for _, pod := range pods {
		newPodNode := createPodNode(pod.Name, uniqueIDs[pod.Name], numConnections[pod.Name])
		nodes = append(nodes, newPodNode)
	}
	return nodes
}

func createPodEdges(pods []models.Pod, uniqueIDs map[string]int) []Edge {
	edges := []Edge{}
	for _, pod := range pods {
		srcID := uniqueIDs[pod.Name]
		for _, dstPod := range pod.Interacts {
			destID := uniqueIDs[dstPod.Name]
			edges = append(edges, createPodEdge(srcID, destID, int(dstPod.Count)))
		}
	}
	return edges
}

func createPodNode(podName string, podID int, podConnections int) Node {
	return Node{
		ID:    podID,
		Label: podName,
		Title: "pods",
		Value: podConnections,
		Group: 1,
	}
}

func createPodEdge(fromID int, toID int, count int) Edge {
	return Edge{
		From:  fromID,
		To:    toID,
		Title: strconv.Itoa(count) + " times communicated",
	}
}
