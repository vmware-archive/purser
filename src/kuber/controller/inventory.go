package controller

type Pod struct {
	name      string
	node      string
	pvcs      []string
	namespace string
	// change the type for time
	startTime float32
	endTime   float32
	labels    []string
}

type Node struct {
	name         string
	instanceType string
}
