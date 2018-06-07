package main

import (
	"fmt"
	"os"
)

func main() {
	inputs := os.Args[1:]
	inputs = inputs[1:]
	if len(inputs) >= 4 && inputs[0] == "get" && inputs[1] == "cost" {
		if inputs[2] == "label" {
			getPodsCostForLabel(inputs[3])
		} else if inputs[2] == "pod" {
			fmt.Println("Work In Progress...")
		} else if inputs[2] == "node" {
			fmt.Println("Work In Progress...")
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("Try one of the following commands...\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost label <key=val>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost pod <pod name>\n")
	fmt.Printf("kubectl --kubeconfig=<absolute path to config> plugin kuber get cost node <node name>\n")
}
