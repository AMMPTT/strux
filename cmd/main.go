package main

import (
    "fmt"
    "time"
    "go-server/internal/ecs"
    "go-server/pkg/components"
)

func main() {
    world := ecs.NewWorld()
    
    // Create breathing system
    breathingSystem := ecs.NewBreathingSystem(world)
    world.AddSystem(breathingSystem)
    
    // Subscribe to breath events
    world.EventManager.Subscribe("EntityBreathed", func(data interface{}) {
        if event, ok := data.(ecs.BreathEvent); ok {
            state := "exhaling"
            if event.State == components.Inhale {
                state = "inhaling"
            }
            fmt.Printf("Entity %d is %s (Volume: %.2f)\n", 
                event.Entity, state, event.Volume)
        }
    })
    
    // Create entity with components
    entity := world.CreateEntity()
    world.AddComponent(entity, &components.Lung{
        Capacity: 1.0,
        Volume:   0.0,
        State:    components.Inhale,
    })
    world.AddComponent(entity, &components.Mouth{
        IsOpen: true,
    })
    
    // Simulation loop
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    fmt.Println("Starting breathing simulation...")
    for i := 0; i < 50; i++ {
        <-ticker.C
        world.Update(0.1)
    }
}