// internal/ecs/component_storage.go

package ecs

import (
    "sync"
    "github.com/AMMPTT/strux/pkg/components"
)

type ComponentArray struct {
    Data        []components.ComponentData
    SparseIndex map[Entity]int
    Size        int
    sync.RWMutex
}


func NewComponentArray() *ComponentArray {
    return &ComponentArray{
        Data:        make([]components.ComponentData, 0),
        SparseIndex: make(map[Entity]int),
    }
}

func (ca *ComponentArray) Add(entity Entity, component components.ComponentData) {
    ca.Lock()
    defer ca.Unlock()

    if index, exists := ca.SparseIndex[entity]; exists {
        ca.Data[index] = component
    } else {
        ca.Data = append(ca.Data, component)
        ca.SparseIndex[entity] = ca.Size
        ca.Size++
    }
}

func (ca *ComponentArray) Remove(entity Entity) {
    ca.Lock()
    defer ca.Unlock()

    if index, exists := ca.SparseIndex[entity]; exists {
        lastIndex := ca.Size - 1
        ca.Data[index] = ca.Data[lastIndex]
        ca.Data = ca.Data[:lastIndex]
        delete(ca.SparseIndex, entity)
        ca.Size--
    }
}

func (ca *ComponentArray) Get(entity Entity) (components.ComponentData, bool) {
    ca.RLock()
    defer ca.RUnlock()

    if index, exists := ca.SparseIndex[entity]; exists {
        return ca.Data[index], true
    }
    return nil, false
}

func (ca *ComponentArray) GetAll() []components.ComponentData {
    ca.RLock()
    defer ca.RUnlock()

    return ca.Data
}