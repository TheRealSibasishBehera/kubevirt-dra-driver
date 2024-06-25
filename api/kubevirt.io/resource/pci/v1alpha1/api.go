package v1alpha1

import nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"

const (
	GroupName = "pci.resource.kubevirt.io"
	Version   = "v1alpha1"

	PciClaimParametersKind = "PciClaimParameters"
)

func DefaultDeviceClassParametersSpec() *DeviceClassParametersSpec {
	return &DeviceClassParametersSpec{
		DeviceSelector: []DeviceSelector{
			{
				Type:              nascrd.PciDeviceType,
				ResourceName:      "*",
				PCIVendorSelector: "*",
			},
		},
	}
}

func DefaultPciClaimParametersSpec() *PciClaimParametersSpec {
	return &PciClaimParametersSpec{
		DeviceName: "*",
		Count:      1,
	}
}
