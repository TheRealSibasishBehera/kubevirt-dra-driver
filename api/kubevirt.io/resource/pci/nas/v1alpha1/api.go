package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupName = "nas.pci.resource.kubevirt.io"
	Version   = "v1alpha1"

	PciDeviceType     = "pci"
	UnknownDeviceType = "unknown"

	NodeAllocationStateStatusReady    = "Ready"
	NodeAllocationStateStatusNotReady = "NotReady"
)

type NodeAllocationStateConfig struct {
	Name      string
	Namespace string
	Owner     *metav1.OwnerReference
}

func NewNodeAllocationState(config *NodeAllocationStateConfig) *NodeAllocationState {
	nascrd := &NodeAllocationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	if config.Owner != nil {
		nascrd.OwnerReferences = []metav1.OwnerReference{*config.Owner}
	}

	return nascrd
}
