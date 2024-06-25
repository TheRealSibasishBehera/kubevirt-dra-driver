package main

import (
	"fmt"
	"sync"

	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
)

type AllocatableDevices map[string]*AllocatableDeviceInfo
type PreparedClaims map[string]*PreparedDevices

type PreparedPcis struct {
	Devices []*PCIDevice
}

type PreparedDevices struct {
	Pci *PreparedPcis
}

func (d PreparedDevices) Type() string {
	if d.Pci != nil {
		return nascrd.PciDeviceType
	}
	return nascrd.UnknownDeviceType
}

type AllocatableDeviceInfo struct {
	*PCIDevice
}

type DeviceState struct {
	sync.Mutex
	cdi           *CDIHandler
	allocatable   AllocatableDevices
	prepared      PreparedClaims
}

func NewDeviceState(config *Config, possibleDevices AllocatableDevices) (*DeviceState, error) {
	cdi, err := NewCDIHandler(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create CDI handler: %v", err)
	}

	err = cdi.CreateCommonSpecFile()
	if err != nil {
		return nil, fmt.Errorf("unable to create CDI spec file for common edits: %v", err)
	}

	state := &DeviceState{
		cdi:           cdi,
		allocatable:   possibleDevices,
		prepared:      make(PreparedClaims),
	}

	err = state.syncPreparedDevicesFromCRDSpec(&config.nascr.Spec)
	if err != nil {
		return nil, fmt.Errorf("unable to sync prepared devices from CRD: %v", err)
	}

	return state, nil
}

func (s *DeviceState) Prepare(claimUID string, allocation nascrd.AllocatedDevices) ([]string, error) {
	s.Lock()
	defer s.Unlock()

	if s.prepared[claimUID] != nil {
		cdiDevices, err := s.cdi.GetClaimDevices(claimUID, s.prepared[claimUID])
		if err != nil {
			return nil, fmt.Errorf("unable to get CDI devices names: %v", err)
		}
		return cdiDevices, nil
	}

	prepared := &PreparedDevices{}

	var err error
	switch allocation.Type() {
	case nascrd.PciDeviceType:
		prepared.Pci, err = s.preparePcis(claimUID, allocation.Pci)
	default:
		err = fmt.Errorf("unknown device type: %v", allocation.Type())
	}
	if err != nil {
		return nil, fmt.Errorf("allocation failed: %v", err)
	}

	err = s.cdi.CreateClaimSpecFile(claimUID, prepared)
	if err != nil {
		return nil, fmt.Errorf("unable to create CDI spec file for claim: %v", err)
	}

	s.prepared[claimUID] = prepared

	cdiDevices, err := s.cdi.GetClaimDevices(claimUID, s.prepared[claimUID])
	if err != nil {
		return nil, fmt.Errorf("unable to get CDI devices names: %v", err)
	}
	return cdiDevices, nil
}

func (s *DeviceState) Unprepare(claimUID string) error {
	s.Lock()
	defer s.Unlock()

	if s.prepared[claimUID] == nil {
		return nil
	}

	switch s.prepared[claimUID].Type() {
	case nascrd.PciDeviceType:
		err := s.unpreparePcis(claimUID, s.prepared[claimUID])
		if err != nil {
			return fmt.Errorf("unprepare failed: %v", err)
		}
	default:
		return fmt.Errorf("unknown device type: %v", s.prepared[claimUID].Type())
	}

	err := s.cdi.DeleteClaimSpecFile(claimUID)
	if err != nil {
		return fmt.Errorf("unable to delete CDI spec file for claim: %v", err)
	}

	delete(s.prepared, claimUID)

	return nil
}

func (s *DeviceState) GetUpdatedSpec(inspec *nascrd.NodeAllocationStateSpec) (*nascrd.NodeAllocationStateSpec, error) {
	s.Lock()
	defer s.Unlock()

	outspec := inspec.DeepCopy()
	err := s.syncAllocatableDevicesToCRDSpec(outspec)
	if err != nil {
		return nil, fmt.Errorf("synching allocatable devices to CR spec: %v", err)
	}

	err = s.syncPreparedDevicesToCRDSpec(outspec)
	if err != nil {
		return nil, fmt.Errorf("synching prepared devices to CR spec: %v", err)
	}

	return outspec, nil
}

func (s *DeviceState) preparePcis(claimUID string, allocated *nascrd.AllocatedPcis) (*PreparedPcis, error) {
	prepared := &PreparedPcis{}

	for _, device := range allocated.Devices {
		pciInfo := s.allocatable[device.UUID].PCIDevice

		if _, exists := s.allocatable[device.UUID]; !exists {
			return nil, fmt.Errorf("requested PCI does not exist: %v", device.UUID)
		}

		prepared.Devices = append(prepared.Devices, pciInfo)
	}

	return prepared, nil
}

func (s *DeviceState) unpreparePcis(claimUID string, devices *PreparedDevices) error {
	return nil
}

func (s *DeviceState) syncAllocatableDevicesToCRDSpec(spec *nascrd.NodeAllocationStateSpec) error {
	pcis := make(map[string]nascrd.AllocatableDevice)
	for _, device := range s.allocatable {
		pcis[device.uuid] = nascrd.AllocatableDevice{
			Pci: &nascrd.AllocatablePci{
				UUID:         device.uuid,
				PciAddress:   device.pciAddress,
				ResourceName: device.resourceName,
			},
		}
	}

	var allocatable []nascrd.AllocatableDevice
	for _, device := range pcis {
		allocatable = append(allocatable, device)
	}

	spec.AllocatableDevices = allocatable

	return nil
}

func (s *DeviceState) syncPreparedDevicesFromCRDSpec(spec *nascrd.NodeAllocationStateSpec) error {
	pcis := s.allocatable

	prepared := make(PreparedClaims)
	for claim, devices := range spec.PreparedClaims {
		switch devices.Type() {
		case nascrd.PciDeviceType:
			prepared[claim] = &PreparedDevices{Pci: &PreparedPcis{}}
			for _, d := range devices.Pci.Devices {
				prepared[claim].Pci.Devices = append(prepared[claim].Pci.Devices, pcis[d.UUID].PCIDevice)
			}
		default:
			return fmt.Errorf("unknown device type: %v", devices.Type())
		}
	}

	s.prepared = prepared

	return nil
}

func (s *DeviceState) syncPreparedDevicesToCRDSpec(spec *nascrd.NodeAllocationStateSpec) error {
	outcas := make(map[string]nascrd.PreparedDevices)
	for claim, devices := range s.prepared {
		var prepared nascrd.PreparedDevices
		switch devices.Type() {
		case nascrd.PciDeviceType:
			prepared.Pci = &nascrd.PreparedPcis{}
			for _, device := range devices.Pci.Devices {
				outdevice := nascrd.PreparedPci{
					UUID: device.uuid,
				}
				prepared.Pci.Devices = append(prepared.Pci.Devices, outdevice)
			}
		default:
			return fmt.Errorf("unknown device type: %v", devices.Type())
		}
		outcas[claim] = prepared
	}

	spec.PreparedClaims = outcas

	return nil
}
