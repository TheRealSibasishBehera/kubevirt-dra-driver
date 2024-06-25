package main

import (
	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
	"sync"
)

type PerNodeAllocatedClaims struct {
	sync.RWMutex
	allocations map[string]map[string]nascrd.AllocatedDevices
}

func NewPerNodeAllocatedClaims() *PerNodeAllocatedClaims {
	return &PerNodeAllocatedClaims{
		allocations: make(map[string]map[string]nascrd.AllocatedDevices),
	}
}

func (p *PerNodeAllocatedClaims) Exists(claimUID, node string) bool {
	p.RLock()
	defer p.RUnlock()

	_, exists := p.allocations[claimUID]
	if !exists {
		return false
	}

	_, exists = p.allocations[claimUID][node]
	return exists
}

func (p *PerNodeAllocatedClaims) Get(claimUID, node string) nascrd.AllocatedDevices {
	p.RLock()
	defer p.RUnlock()

	if !p.Exists(claimUID, node) {
		return nascrd.AllocatedDevices{}
	}
	return p.allocations[claimUID][node]
}

func (p *PerNodeAllocatedClaims) VisitNode(node string, visitor func(claimUID string, allocation nascrd.AllocatedDevices)) {
	p.RLock()
	for claimUID := range p.allocations {
		if allocation, exists := p.allocations[claimUID][node]; exists {
			p.RUnlock()
			visitor(claimUID, allocation)
			p.RLock()
		}
	}
	p.RUnlock()
}

func (p *PerNodeAllocatedClaims) Visit(visitor func(claimUID, node string, allocation nascrd.AllocatedDevices)) {
	p.RLock()
	for claimUID := range p.allocations {
		for node, allocation := range p.allocations[claimUID] {
			p.RUnlock()
			visitor(claimUID, node, allocation)
			p.RLock()
		}
	}
	p.RUnlock()
}

func (p *PerNodeAllocatedClaims) Set(claimUID, node string, devices nascrd.AllocatedDevices) {
	p.Lock()
	defer p.Unlock()

	_, exists := p.allocations[claimUID]
	if !exists {
		p.allocations[claimUID] = make(map[string]nascrd.AllocatedDevices)
	}

	p.allocations[claimUID][node] = devices
}

func (p *PerNodeAllocatedClaims) RemoveNode(claimUID, node string) {
	p.Lock()
	defer p.Unlock()

	_, exists := p.allocations[claimUID]
	if !exists {
		return
	}

	delete(p.allocations[claimUID], node)
}

func (p *PerNodeAllocatedClaims) Remove(claimUID string) {
	p.Lock()
	defer p.Unlock()

	delete(p.allocations, claimUID)
}
