package main

import (
	"fmt"

	pcicrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/api/resource/v1alpha2"
	"k8s.io/dynamic-resource-allocation/controller"

	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
)

type pcidriver struct {
	PendingAllocatedClaims *PerNodeAllocatedClaims
}

func NewPciDriver() *pcidriver {
	return &pcidriver{
		PendingAllocatedClaims: NewPerNodeAllocatedClaims(),
	}
}

func (p *pcidriver) ValidateClaimParameters(claimParams *pcicrd.PciClaimParametersSpec) error {
	if claimParams.DeviceName != "devices.kubevirt.io/nvme" {
		return fmt.Errorf("unsupported pci device type: %s", claimParams.DeviceName)
	}
	return nil
}

func (p *pcidriver) Allocate(crd *nascrd.NodeAllocationState, claim *resourcev1.ResourceClaim, claimParams *pcicrd.PciClaimParametersSpec, class *resourcev1.ResourceClass, classParams *pcicrd.DeviceClassParametersSpec, selectedNode string) (OnSuccessCallback, error) {
	claimUID := string(claim.UID)

	if !p.PendingAllocatedClaims.Exists(claimUID, selectedNode) {
		return nil, fmt.Errorf("no allocations generated for claim '%v' on node '%v' yet", claim.UID, selectedNode)
	}

	crd.Spec.AllocatedClaims[claimUID] = p.PendingAllocatedClaims.Get(claimUID, selectedNode)
	onSuccess := func() {
		p.PendingAllocatedClaims.Remove(claimUID)
	}

	return onSuccess, nil
}

func (p *pcidriver) Deallocate(crd *nascrd.NodeAllocationState, claim *resourcev1.ResourceClaim) error {
	claimUID := string(claim.UID)
	p.PendingAllocatedClaims.Remove(claimUID)
	return nil
}

func (p *pcidriver) UnsuitableNode(crd *nascrd.NodeAllocationState, pod *corev1.Pod, pcicas []*controller.ClaimAllocation, allcas []*controller.ClaimAllocation, potentialNode string) error {

	p.PendingAllocatedClaims.VisitNode(potentialNode, func(claimUID string, allocation nascrd.AllocatedDevices) {
		if _, exists := crd.Spec.AllocatedClaims[claimUID]; exists {
			p.PendingAllocatedClaims.Remove(claimUID)
		} else {
			crd.Spec.AllocatedClaims[claimUID] = allocation
		}
	})

	allocated := p.allocate(crd, pod, pcicas, allcas, potentialNode)
	for _, ca := range pcicas {
		claimUID := string(ca.Claim.UID)
		claimParams, _ := ca.ClaimParameters.(*pcicrd.PciClaimParametersSpec)

		//TODO :Remove count
		if claimParams.Count != len(allocated[claimUID]) {
			for _, ca := range allcas {
				ca.UnsuitableNodes = append(ca.UnsuitableNodes, potentialNode)
			}
			return nil
		}

		var devices []nascrd.AllocatedPci
		for _, pci := range allocated[claimUID] {
			device := nascrd.AllocatedPci{
				UUID: pci,
			}
			devices = append(devices, device)
		}

		allocatedDevices := nascrd.AllocatedDevices{
			Pci: &nascrd.AllocatedPcis{
				Devices: devices,
			},
		}

		p.PendingAllocatedClaims.Set(claimUID, potentialNode, allocatedDevices)
	}

	return nil
}

func (p *pcidriver) allocate(crd *nascrd.NodeAllocationState, pod *corev1.Pod, pcicas []*controller.ClaimAllocation, allcas []*controller.ClaimAllocation, node string) map[string][]string {

	available := make(map[string]*nascrd.AllocatablePci)

	for _, device := range crd.Spec.AllocatableDevices {
		switch device.Type() {
		case nascrd.PciDeviceType:
			available[device.Pci.UUID] = device.Pci
		default:
			// skip other devices
		}
	}

	for _, allocation := range crd.Spec.AllocatedClaims {
		switch allocation.Type() {
		case nascrd.PciDeviceType:
			for _, device := range allocation.Pci.Devices {
				delete(available, device.UUID)
			}
		default:
		}
	}

	allocated := make(map[string][]string)

	for _, ca := range pcicas {
		claimUID := string(ca.Claim.UID)

		if _, exists := crd.Spec.AllocatedClaims[claimUID]; exists {
			devices := crd.Spec.AllocatedClaims[claimUID].Pci.Devices
			for _, device := range devices {
				allocated[claimUID] = append(allocated[claimUID], device.UUID)
			}
			continue
		}

		claimParams, _ := ca.ClaimParameters.(*pcicrd.PciClaimParametersSpec)
		var devices []string

		for i := 0; i < claimParams.Count; i++ {
			for uuid, device := range available {

				//check the device type is the one requested
				if device.ResourceName == claimParams.DeviceName {
					devices = append(devices, device.UUID)
					delete(available, uuid)
					break
				}
			}
		}
		allocated[claimUID] = devices
	}

	return allocated
}
