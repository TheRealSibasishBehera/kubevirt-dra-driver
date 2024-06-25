package main

import (
	"strings"
)

func enumerateAllPossibleDevices(supportedAllocatedDevices AllocatableDevices) (AllocatableDevices) {
	supportedPCIDeviceMap := make(map[string]string)
	for _, device := range supportedAllocatedDevices {
		supportedPCIDeviceMap[strings.ToLower(device.PCIDevice.vendorSelector)] = device.PCIDevice.resourceName
	}

	allDevices := make(AllocatableDevices)
	discoveredDevices := DiscoverPermittedHostPCIDevices(supportedPCIDeviceMap)
	for _, device := range supportedAllocatedDevices {
		vendorSelector := device.PCIDevice.vendorSelector
		if devices, supported := discoveredDevices[vendorSelector]; supported {
			for _, discoveredDevice := range devices {
				newDevice := &AllocatableDeviceInfo{
					PCIDevice: discoveredDevice,
				}
				allDevices[newDevice.PCIDevice.uuid] = newDevice
			}
		}
	}

	return allDevices
}
