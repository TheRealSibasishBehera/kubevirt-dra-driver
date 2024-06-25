package main

import (
	"sync"
)

type PerNodeMutex struct {
	sync.Mutex
	submutex map[string]*sync.Mutex
}

func NewPerNodeMutex() *PerNodeMutex {
	return &PerNodeMutex{
		submutex: make(map[string]*sync.Mutex),
	}
}

func (pnm *PerNodeMutex) Get(node string) *sync.Mutex {
	pnm.Mutex.Lock()
	defer pnm.Mutex.Unlock()
	if pnm.submutex[node] == nil {
		pnm.submutex[node] = &sync.Mutex{}
	}
	return pnm.submutex[node]
}
