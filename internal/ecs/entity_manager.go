package ecs

import (
	"fmt"
	"reflect"
	"sync"
	"go-server/pkg/components"
	
)

type EntityManager struct {
    entities       map[Entity]bool
    components     map[Entity]map[reflect.Type]components.ComponentData
    archetypes     []*Archetype
    componentPools map[reflect.Type]*ComponentPool
    nextEntityID   Entity
    mu             sync.RWMutex
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities:       make(map[Entity]bool),
		components:     make(map[Entity]map[reflect.Type]Component),
		archetypes:     make([]*Archetype, 0),
		componentPools: make(map[reflect.Type]*ComponentPool),
	}
}

func (em *EntityManager) CreateEntity() Entity {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	id := em.nextEntityID
	em.entities[id] = true
	em.components[id] = make(map[reflect.Type]Component)
	em.nextEntityID++
	return id
}

func (em *EntityManager) DestroyEntity(entity Entity) {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	if !em.entities[entity] {
		return // Entity doesn't exist, nothing to do
	}
	
	delete(em.entities, entity)
	for componentType, component := range em.components[entity] {
		if pool, exists := em.componentPools[componentType]; exists {
			pool.Return(component)
		}
	}
	delete(em.components, entity)
	
	for _, archetype := range em.archetypes {
		for i, e := range archetype.entities {
			if e == entity {
				lastIdx := len(archetype.entities) - 1
				archetype.entities[i] = archetype.entities[lastIdx]
				archetype.entities = archetype.entities[:lastIdx]
				for j := range archetype.components {
					archetype.components[j][i] = archetype.components[j][lastIdx]
					archetype.components[j] = archetype.components[j][:lastIdx]
				}
				break
			}
		}
	}
}

func (em *EntityManager) UpdateComponent(entity Entity, component Component) {
    em.mu.Lock()
    defer em.mu.Unlock()

    if !em.entities[entity] {
        panic(fmt.Sprintf("Entity %d does not exist", entity))
    }

    componentType := reflect.TypeOf(component)
    if _, exists := em.components[entity][componentType]; !exists {
        panic(fmt.Sprintf("Component of type %v does not exist for entity %d", componentType, entity))
    }

    em.components[entity][componentType] = component

    // Update component in the appropriate archetype
    for _, archetype := range em.archetypes {
        for i, e := range archetype.entities {
            if e == entity {
                for j, ct := range archetype.componentTypes {
                    if ct == componentType {
                        archetype.components[j][i] = component
                        break
                    }
                }
                break
            }
        }
    }

    fmt.Printf("Updated Component!... %d %+v\n", entity, component)
}

func (em *EntityManager) AddComponent(entity Entity, component Component) {
	em.mu.Lock()
	defer em.mu.Unlock()
	fmt.Println("Adding Component!...", entity, component)

	
	if !em.entities[entity] {
		panic(fmt.Sprintf("Entity %d does not exist", entity))
	}
	
	componentType := reflect.TypeOf(component)
	if pool, exists := em.componentPools[componentType]; exists {
		component = pool.Get()
	}
	em.components[entity][componentType] = component
	
	// Update archetypes
	componentTypes := make([]reflect.Type, 0, len(em.components[entity]))
	for ct := range em.components[entity] {
		componentTypes = append(componentTypes, ct)
	}
	
	var targetArchetype *Archetype
	for _, archetype := range em.archetypes {
		if reflect.DeepEqual(archetype.componentTypes, componentTypes) {
			targetArchetype = archetype
			break
		}
	}
	
	if targetArchetype == nil {
		targetArchetype = &Archetype{
			componentTypes: componentTypes,
			entities:       make([]Entity, 0),
			components:     make([][]Component, len(componentTypes)),
		}
		em.archetypes = append(em.archetypes, targetArchetype)
	}
	
	targetArchetype.entities = append(targetArchetype.entities, entity)
	for i, ct := range targetArchetype.componentTypes {
		targetArchetype.components[i] = append(targetArchetype.components[i], em.components[entity][ct])
	}
}

func (em *EntityManager) RemoveComponent(entity Entity, componentType reflect.Type) {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	if !em.entities[entity] {
		return // Entity doesn't exist, nothing to do
	}
	
	if component, exists := em.components[entity][componentType]; exists {
		if pool, exists := em.componentPools[componentType]; exists {
			pool.Return(component)
		}
		delete(em.components[entity], componentType)
		
		// Update archetypes
		for _, archetype := range em.archetypes {
			for i, e := range archetype.entities {
				if e == entity {
					lastIdx := len(archetype.entities) - 1
					archetype.entities[i] = archetype.entities[lastIdx]
					archetype.entities = archetype.entities[:lastIdx]
					for j := range archetype.components {
						archetype.components[j][i] = archetype.components[j][lastIdx]
						archetype.components[j] = archetype.components[j][:lastIdx]
					}
					break
				}
			}
		}
	}
}

func (em *EntityManager) GetComponent(entity Entity, componentType reflect.Type) (Component, bool) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	
	if components, exists := em.components[entity]; exists {
		component, exists := components[componentType]
		return component, exists
	}
	return nil, false
}

func (em *EntityManager) Query(componentTypes ...reflect.Type) []Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()
	
	var result []Entity
	for _, archetype := range em.archetypes {
		if len(archetype.componentTypes) < len(componentTypes) {
			continue
		}
		
		match := true
		for _, queryType := range componentTypes {
			found := false
			for _, archetypeType := range archetype.componentTypes {
				if archetypeType == queryType {
					found = true
					break
				}
			}
			if !found {
				match = false
				break
			}
		}
		
		if match {
			result = append(result, archetype.entities...)
		}
	}
	return result
}

func (em *EntityManager) InitializeComponentPool(componentType reflect.Type) {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	if _, exists := em.componentPools[componentType]; !exists {
		em.componentPools[componentType] = NewComponentPool(componentType)
	}
}