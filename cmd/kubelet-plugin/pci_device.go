package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	pciBasePath = "/sys/bus/pci/devices"
	PCIResourcePrefix   = "PCI_RESOURCE"
)

type PCIDevice struct {
	uuid           string
	vendorSelector string
	resourceName   string
	pciAddress     string
	driver         string
	iommuGroup     string
	numaNode       int
	pciID          string
}

func DiscoverPermittedHostPCIDevices(supportedPCIDeviceMap map[string]string) (map[string][]*PCIDevice) {
	initHandler()

	iommuToPCIMap := make(map[string]string)

	pciDevicesMap := make(map[string][]*PCIDevice)
	err := filepath.Walk(pciBasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		pciID, err := Handler.GetDevicePCIID(pciBasePath, info.Name())
		if err != nil {
            log.Printf("Failed to get vendor:device ID for device: %s, error: %v", info.Name(), err)
			return nil 
		}

		if _, supported := supportedPCIDeviceMap[pciID]; supported {
			driver, err := Handler.GetDeviceDriver(pciBasePath, info.Name())
			if err != nil || driver != "vfio-pci" {
				log.Printf("Driver error: %v", err)
				return nil 
			}

			pcidev := &PCIDevice{
				uuid:           uuid.New().String(),
				pciID:          pciID,
				pciAddress:     info.Name(),
				vendorSelector: pciID,
				resourceName:   supportedPCIDeviceMap[pciID],
			}

			iommuGroup, err := Handler.GetDeviceIOMMUGroup(pciBasePath, info.Name())
			if err != nil {
				log.Printf("IOMMU group error: %v", err)
                return nil
			}
			pcidev.iommuGroup = iommuGroup
			pcidev.driver = driver
			pcidev.numaNode = Handler.GetDeviceNumaNode(pciBasePath, info.Name())

			iommuToPCIMap[pcidev.iommuGroup] = pcidev.pciAddress

			pciDevicesMap[pciID] = append(pciDevicesMap[pciID], pcidev)
		}
		return nil
	})
	if err != nil {
        log.Printf("Failed to discover host devices, error: %v", err)
	}

	return pciDevicesMap
}

// MockDiscoverPermittedHostPCIDevices returns predefined data for testing
func MockDiscoverPermittedHostPCIDevices(supportedPCIDeviceMap map[string]string) (map[string][]*PCIDevice, map[string]string) {
	log.Printf("enter  MockDiscoverPermittedHostPCIDevices")
	pciDevicesMap := make(map[string][]*PCIDevice)
	iommuToPCIMap := make(map[string]string)
	mapLength := len(supportedPCIDeviceMap)
	log.Println("The length of the map is:", mapLength)

	for vendorSelector, resourceName := range supportedPCIDeviceMap {
		pcidev1 := &PCIDevice{
			uuid:           uuid.New().String(), 
			vendorSelector: vendorSelector,
			resourceName:   resourceName,
			pciAddress:     "0000:00:1d." + vendorSelector, 
			driver:         "vfio-pci",                     
			iommuGroup:     "20" + vendorSelector,          
			numaNode:       0,                              
			pciID:          vendorSelector,
		}

		pcidev2 := &PCIDevice{
			uuid:           uuid.New().String(), 
			vendorSelector: vendorSelector,
			resourceName:   resourceName,
			pciAddress:     "0000:00:1e." + vendorSelector, 
			driver:         "vfio-pci",                     
			iommuGroup:     "30" + vendorSelector,          
			numaNode:       0,                              
			pciID:          vendorSelector,
		}
		iommuToPCIMap["20"+vendorSelector] = "0000:00:1d." + vendorSelector 
		iommuToPCIMap["30"+vendorSelector] = "0000:00:1e." + vendorSelector 

		pciDevicesMap[vendorSelector] = append(pciDevicesMap[vendorSelector], pcidev1, pcidev2)

	}

	return pciDevicesMap, iommuToPCIMap
}
