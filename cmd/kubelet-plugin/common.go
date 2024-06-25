package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"log"
)

type DeviceHandler interface {
	GetDeviceIOMMUGroup(basepath string, pciAddress string) (string, error)
	GetDeviceDriver(basepath string, pciAddress string) (string, error)
	GetDeviceNumaNode(basepath string, pciAddress string) (numaNode int)
	GetDevicePCIID(basepath string, pciAddress string) (string, error)
}

type DeviceUtilsHandler struct{}

var Handler DeviceHandler

// getDeviceIOMMUGroup gets devices iommu_group
// e.g. /sys/bus/pci/devices/0000\:65\:00.0/iommu_group -> ../../../../../kernel/iommu_groups/45
func (h *DeviceUtilsHandler) GetDeviceIOMMUGroup(basepath string, pciAddress string) (string, error) {
	iommuLink := filepath.Join(basepath, pciAddress, "iommu_group")
	iommuPath, err := os.Readlink(iommuLink)
	if err != nil {
		log.Printf("failed to read iommu_group link %s for device %s", iommuLink, pciAddress)
		return "", err
	}
	_, iommuGroup := filepath.Split(iommuPath)
	return iommuGroup, nil
}

// gets device driver
func (h *DeviceUtilsHandler) GetDeviceDriver(basepath string, pciAddress string) (string, error) {
	driverLink := filepath.Join(basepath, pciAddress, "driver")
	driverPath, err := os.Readlink(driverLink)
	if err != nil {
		log.Printf("failed to read driver link %s for device %s", driverLink, pciAddress)
		return "", err
	}
	_, driver := filepath.Split(driverPath)
	return driver, nil
}

func (h *DeviceUtilsHandler) GetDeviceNumaNode(basepath string, pciAddress string) (numaNode int) {
	numaNode = -1
	numaNodePath := filepath.Join(basepath, pciAddress, "numa_node")
	// #nosec No risk for path injection. Reading static path of NUMA node info
	numaNodeStr, err := os.ReadFile(numaNodePath)
	if err != nil {
		log.Printf("failed to read numa_node %s for device %s", numaNodePath, pciAddress)
		return
	}
	numaNodeStr = bytes.TrimSpace(numaNodeStr)
	numaNode, err = strconv.Atoi(string(numaNodeStr))
	if err != nil {
		log.Printf("failed to convert numa node value %v of device %s", numaNodeStr, pciAddress)
		return
	}
	return
}

func (h *DeviceUtilsHandler) GetDevicePCIID(basepath string, pciAddress string) (string, error) {
	log.Printf("GetDevicePCIID: basepath: %s, pciAddress: %s", basepath, pciAddress)
	// #nosec No risk for path injection. Reading static path of PCI data
	file, err := os.Open(filepath.Join(basepath, pciAddress, "uevent"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PCI_ID") {
			equal := strings.Index(line, "=")
			value := strings.TrimSpace(line[equal+1:])
			return strings.ToLower(value), nil
		}
	}
	return "", fmt.Errorf("no pci_id is found")
}

func initHandler() {
	if Handler == nil {
		Handler = &DeviceUtilsHandler{}
	}
}
