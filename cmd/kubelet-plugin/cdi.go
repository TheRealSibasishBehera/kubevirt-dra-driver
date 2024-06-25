package main

import (
	"fmt"

	cdiapi "github.com/container-orchestrated-devices/container-device-interface/pkg/cdi"
	cdispec "github.com/container-orchestrated-devices/container-device-interface/specs-go"

	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
	"kubevirt.io/kubevirt-dra-driver/pkg/util"
)

const (
	cdiVendor = "k8s." + DriverName
	cdiClass  = "pci"
	cdiKind   = cdiVendor + "/" + cdiClass
	cdiCommonDeviceName = "common"
)

type CDIHandler struct {
	registry cdiapi.Registry
}

func NewCDIHandler(config *Config) (*CDIHandler, error) {
	registry := cdiapi.GetRegistry(
		cdiapi.WithSpecDirs(config.flags.cdiRoot),
	)

	err := registry.Refresh()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh the CDI registry: %v", err)
	}

	handler := &CDIHandler{
		registry: registry,
	}

	return handler, nil
}

func (cdi *CDIHandler) GetDevice(device string) *cdiapi.Device {
	return cdi.registry.DeviceDB().GetDevice(device)
}

func (cdi *CDIHandler) CreateCommonSpecFile() error {
	spec := &cdispec.Spec{
		Kind: cdiKind,
		Devices: []cdispec.Device{
			{
				Name: cdiCommonDeviceName,
				ContainerEdits: cdispec.ContainerEdits{
					Env: []string{
						fmt.Sprintf("DRA_RESOURCE_DRIVER_NAME=%s", DriverName),
					},
				},
			},
		},
	}

	minVersion, err := cdiapi.MinimumRequiredVersion(spec)
	if err != nil {
		return fmt.Errorf("failed to get minimum required CDI spec version: %v", err)
	}
	spec.Version = minVersion

	specName, err := cdiapi.GenerateNameForTransientSpec(spec, cdiCommonDeviceName)
	if err != nil {
		return fmt.Errorf("failed to generate Spec name: %w", err)
	}

	return cdi.registry.SpecDB().WriteSpec(spec, specName)
}

func (cdi *CDIHandler) CreateClaimSpecFile(claimUID string, devices *PreparedDevices) error {
	specName := cdiapi.GenerateTransientSpecName(cdiVendor, cdiClass, claimUID)

	spec := &cdispec.Spec{
		Kind:    cdiKind,
		Devices: []cdispec.Device{},
	}
	switch devices.Type() {
	case nascrd.PciDeviceType:
		for _, device := range devices.Pci.Devices {

			resouceNameEnvVar := util.ResourceNameToEnvVar(PCIResourcePrefix, device.resourceName)

			cdiDevice := cdispec.Device{
				Name: device.uuid,
				ContainerEdits: cdispec.ContainerEdits{
					Env: []string{
						fmt.Sprintf(resouceNameEnvVar+"=%s", device.pciAddress),
					},
				},
			}
			spec.Devices = append(spec.Devices, cdiDevice)
		}
	default:
		return fmt.Errorf("unknown device type: %v", devices.Type())
	}

	minVersion, err := cdiapi.MinimumRequiredVersion(spec)
	if err != nil {
		return fmt.Errorf("failed to get minimum required CDI spec version: %v", err)
	}
	spec.Version = minVersion

	return cdi.registry.SpecDB().WriteSpec(spec, specName)
}

func (cdi *CDIHandler) DeleteClaimSpecFile(claimUID string) error {
	specName := cdiapi.GenerateTransientSpecName(cdiVendor, cdiClass, claimUID)
	return cdi.registry.SpecDB().RemoveSpec(specName)
}

func (cdi *CDIHandler) GetClaimDevices(claimUID string, devices *PreparedDevices) ([]string, error) {
	cdiDevices := []string{
		cdiapi.QualifiedName(cdiVendor, cdiClass, cdiCommonDeviceName),
	}

	switch devices.Type() {
	case nascrd.PciDeviceType:
		for _, device := range devices.Pci.Devices {
			cdiDevice := cdiapi.QualifiedName(cdiVendor, cdiClass, device.uuid)
			cdiDevices = append(cdiDevices, cdiDevice)
		}
	default:
		return nil, fmt.Errorf("unknown device type: %v", devices.Type())
	}

	return cdiDevices, nil
}
