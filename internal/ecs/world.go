package ecs

import (
    "encoding/json"
    "reflect"
    "sync"
    "fmt"
    "github.com/AMMPTT/strux/pkg/components"
)

type World struct {
    entities      map[Entity]bool
    components    map[reflect.Type]*ComponentArray
    systems       []System
    EventManager  *EventManager  // Changed to uppercase to export
    mu            sync.RWMutex
}

func NewWorld() *World {
    return &World{
        entities:     make(map[Entity]bool),
        components:   make(map[reflect.Type]*ComponentArray),
        systems:      make([]System, 0),
        EventManager: NewEventManager(),
    }
}


func (w *World) AddSystem(system System) {
    w.systems = append(w.systems, system)
    fmt.Println("Adding System!...", system)
}

func (w *World) Update(dt float32) {
    var wg sync.WaitGroup
    for _, system := range w.systems {
        wg.Add(1)
        go func(s System) {
            defer wg.Done()
            s.Update(dt)
        }(system)
    }
    wg.Wait()
}

func (w *World) CreateEntity() Entity {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    id := Entity(len(w.entities))
    w.entities[id] = true
    return id
}

func (w *World) AddComponent(entity Entity, component components.ComponentData) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    componentType := reflect.TypeOf(component)
    if w.components[componentType] == nil {
        w.components[componentType] = NewComponentArray()
    }
    w.components[componentType].Add(entity, component)
}

func (w *World) RemoveComponent(entity Entity, componentType reflect.Type) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if compArray, exists := w.components[componentType]; exists {
        compArray.Remove(entity)
    }
}

func (w *World) GetComponent(entity Entity, componentType reflect.Type) (components.ComponentData, bool) {
    w.mu.RLock()
    defer w.mu.RUnlock()
    
    if compArray, exists := w.components[componentType]; exists {
        return compArray.Get(entity)
    }
    return nil, false
}

func (w *World) SaveState() ([]byte, error) {
    w.mu.RLock()
    defer w.mu.RUnlock()
    
    state := struct {
        Entities   map[Entity]bool
        Components map[string][]components.ComponentData
    }{
        Entities:   w.entities,
        Components: make(map[string][]components.ComponentData),
    }
    
    for compType, compArray := range w.components {
        state.Components[compType.String()] = compArray.GetAll()
    }
    
    return json.Marshal(state)
}

func (w *World) LoadState(data []byte) error {
    var state struct {
        Entities   map[Entity]bool
        Components map[string][]components.ComponentData
    }
    
    if err := json.Unmarshal(data, &state); err != nil {
        return err
    }
    
    w.mu.Lock()
    defer w.mu.Unlock()
    
    w.entities = state.Entities
    w.components = make(map[reflect.Type]*ComponentArray)
    
    for _, comps := range state.Components {
        if len(comps) == 0 {
            continue
        }
        compType := reflect.TypeOf(comps[0])
        compArray := NewComponentArray()
        for _, comp := range comps {
            compArray.Add(Entity(len(compArray.Data)), comp)
        }
        w.components[compType] = compArray
    }
    
    return nil
}