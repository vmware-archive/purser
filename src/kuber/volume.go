package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolume details
type PersistentVolume struct {
	name         string
	capacityInGB float64
	storageClass string
}

// PersistentVolumeClaim details
type PersistentVolumeClaim struct {
	name                string
	volumeName          string
	requestSizeInGB     float64
	capacityAllotedInGB float64
	storageClass        *string
	pv                  *PersistentVolume
}

func collectPersistentVolumeClaim(claimName string) *PersistentVolumeClaim {
	pvc, err := ClientSetInstance.CoreV1().PersistentVolumeClaims("default").Get(claimName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Persistent Volume Claim %s not found\n", claimName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting persistence volume Claim %s : %v\n", claimName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		var claim PersistentVolumeClaim
		claim.name = pvc.GetObjectMeta().GetName()
		claim.volumeName = pvc.Spec.VolumeName
		claim.storageClass = pvc.Spec.StorageClassName
		request := pvc.Spec.Resources.Requests["storage"].DeepCopy()
		// TODO: consider quantity format(Gi,Mi,G,etc.) into consideration.
		claim.requestSizeInGB = (float64)(request.Value()) / (float64)(1024.0*1024.0*1024.0)
		capacity := pvc.Status.Capacity["storage"].DeepCopy()
		// TODO: consider quantity format(Gi,Mi,G,etc.) into consideration.
		claim.capacityAllotedInGB = (float64)(capacity.Value()) / (float64)(1024.0*1024.0*1024.0)
		fmt.Println(claim)
		return &claim
	}
}

func collectPersistentVolume(volName string) *PersistentVolume {
	pv, err := ClientSetInstance.CoreV1().PersistentVolumes().Get(volName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Persistent Volume %s not found\n", volName)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting persistence volume %s : %v\n", volName, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		var persistentVolume PersistentVolume
		persistentVolume.storageClass = pv.Spec.StorageClassName
		persistentVolume.name = pv.Name
		q := pv.Spec.Capacity["storage"].DeepCopy()
		// TODO: consider quantity format(Gi,Mi,G,etc.) into consideration.
		persistentVolume.capacityInGB = (float64)(q.Value()) / (float64)(1024.0*1024.0*1024.0)
		fmt.Println(persistentVolume)
		return &persistentVolume
	}
}
