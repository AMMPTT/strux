package ecs

import (
	"reflect"
	"sync"
	"go-server/pkg/components"

)

// ComponentPool manages a pool of components of a specific type
type ComponentPool struct {
    componentType reflect.Type
    pool          []components.ComponentData
    mu            sync.Mutex
}

// NewComponentPool creates a new ComponentPool for a given component type
func NewComponentPool(componentType reflect.Type) *ComponentPool {
    return &ComponentPool{
        componentType: componentType,
        pool:          make([]components.ComponentData, 0),
    }
}

// Get retrieves a component from the pool or creates a new one if the pool is empty
func (cp *ComponentPool) Get() Component {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if len(cp.pool) == 0 {
		// Create a new component if the pool is empty
		return reflect.New(cp.componentType).Interface().(Component)
	}

	// Remove and return the last component from the pool
	component := cp.pool[len(cp.pool)-1]
	cp.pool = cp.pool[:len(cp.pool)-1]
	return component
}

// Return puts a component back into the pool
func (cp *ComponentPool) Return(component Component) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Ensure the component is of the correct type
	if reflect.TypeOf(component) != cp.componentType {
		panic("Attempted to return component of wrong type to pool")
	}

	cp.pool = append(cp.pool, component)
}

// Size returns the current number of components in the pool
func (cp *ComponentPool) Size() int {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return len(cp.pool)
}