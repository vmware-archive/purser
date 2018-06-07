package main

// PersistentVolume details
type PersistentVolume struct {
	kind         string
	capacity     float64
	storageClass string
}

// PersistentVolumeClaim details
type PersistentVolumeClaim struct {
	volumeName      string
	requestSize     float64
	capacityAlloted float64
	name            string
	pv              *PersistentVolume
}
