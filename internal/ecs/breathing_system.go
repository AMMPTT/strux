// internal/ecs/breathing_system.go
package ecs

import (
    "reflect"
    "go-server/pkg/components"
)

type BreathingSystem struct {
    world *World
}

func NewBreathingSystem(world *World) *BreathingSystem {
    return &BreathingSystem{world: world}
}

type BreathEvent struct {
    Entity Entity
    State  components.LungState
    Volume float32
}

func (s *BreathingSystem) Update(dt float32) {
    s.world.mu.RLock()
    defer s.world.mu.RUnlock()

    lungArray := s.world.components[reflect.TypeOf(&components.Lung{})]
    mouthArray := s.world.components[reflect.TypeOf(&components.Mouth{})]

    if lungArray == nil || mouthArray == nil {
        return
    }

    lungArray.Lock()
    mouthArray.Lock()
    defer lungArray.Unlock()
    defer mouthArray.Unlock()

    for i := 0; i < lungArray.Size; i++ {
        entityIndex := i
        lung, ok := lungArray.Data[i].(*components.Lung)
        if !ok {
            continue
        }
        mouth, ok := mouthArray.Data[entityIndex].(*components.Mouth)
        if !ok {
            continue
        }

        previousState := lung.State
        previousVolume := lung.Volume

        // Update breathing state
        if lung.State == components.Exhale {
            lung.Volume -= dt * 0.5 // Exhale rate
            if lung.Volume <= 0 {
                lung.Volume = 0
                lung.State = components.Inhale
                mouth.IsOpen = true
            }
        } else {
            lung.Volume += dt * 0.5 // Inhale rate
            if lung.Volume >= lung.Capacity {
                lung.Volume = lung.Capacity
                lung.State = components.Exhale
                mouth.IsOpen = false
            }
        }

        // Emit event if state changed
        if lung.State != previousState || lung.Volume != previousVolume {
            s.world.EventManager.Publish("EntityBreathed", BreathEvent{
                Entity: Entity(entityIndex),
                State:  lung.State,
                Volume: lung.Volume,
            })
        }
    }
}