package main

import (
	"fmt"
	"time"
)

// Component interfaces
type Component interface{}

// MouthComponent represents the state of the mouth
type MouthComponent struct {
	IsOpen bool
}

// LungComponent represents the state of the lungs
type LungComponent struct {
	Capacity     float64
	CurrentVolume float64
}

// BreathingComponent represents the breathing state
type BreathingComponent struct {
	IsInhaling bool
}

// Entity represents a game object
type Entity struct {
	ID         int
	Components map[string]Component
}

// System interface
type System interface {
	Update(entities []*Entity)
}

// BreathingSystem handles the breathing logic
type BreathingSystem struct{}

func (bs *BreathingSystem) Update(entities []*Entity) {
	for _, entity := range entities {
		mouth, hasMouth := entity.Components["Mouth"].(*MouthComponent)
		lung, hasLung := entity.Components["Lung"].(*LungComponent)
		breathing, hasBreathing := entity.Components["Breathing"].(*BreathingComponent)

		if !hasMouth || !hasLung || !hasBreathing {
			continue
		}

		if breathing.IsInhaling {
			if lung.CurrentVolume < lung.Capacity {
				mouth.IsOpen = true
				lung.CurrentVolume += 0.1 // Inhale
				if lung.CurrentVolume >= lung.Capacity {
					breathing.IsInhaling = false
				}
			}
		} else {
			if lung.CurrentVolume > 0 {
				mouth.IsOpen = true
				lung.CurrentVolume -= 0.1 // Exhale
				if lung.CurrentVolume <= 0 {
					breathing.IsInhaling = true
				}
			}
		}

		fmt.Printf("Entity %d - Mouth: %v, Lung: %.2f, Inhaling: %v\n",
			entity.ID, mouth.IsOpen, lung.CurrentVolume, breathing.IsInhaling)
	}
}

func main() {
	// Create an entity
	entity := &Entity{
		ID: 1,
		Components: map[string]Component{
			"Mouth":     &MouthComponent{IsOpen: false},
			"Lung":      &LungComponent{Capacity: 1.0, CurrentVolume: 0},
			"Breathing": &BreathingComponent{IsInhaling: true},
		},
	}

	// Create a breathing system
	breathingSystem := &BreathingSystem{}

	// Simulation loop
	ticker := time.NewTicker(600 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 20; i++ { // Run for 20 ticks
		<-ticker.C
		fmt.Printf("\nTick %d\n", i+1)
		breathingSystem.Update([]*Entity{entity})
	}
}