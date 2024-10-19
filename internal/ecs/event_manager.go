package ecs

import (
    "fmt"
    "sync"
)

type EventManager struct {
    subscribers map[string]map[uint64]func(interface{})
    mu          sync.RWMutex
    nextID      uint64
}

func NewEventManager() *EventManager {
    return &EventManager{
        subscribers: make(map[string]map[uint64]func(interface{})),
        nextID:      1,
    }
}

func (em *EventManager) Subscribe(eventType string, callback func(interface{})) uint64 {
    em.mu.Lock()
    defer em.mu.Unlock()

    if em.subscribers[eventType] == nil {
        em.subscribers[eventType] = make(map[uint64]func(interface{}))
    }

    id := em.nextID
    em.subscribers[eventType][id] = callback
    em.nextID++

    fmt.Printf("Subscribed to %s with ID %d\n", eventType, id)
    return id
}

func (em *EventManager) Unsubscribe(eventType string, id uint64) {
    em.mu.Lock()
    defer em.mu.Unlock()

    if callbacks, exists := em.subscribers[eventType]; exists {
        delete(callbacks, id)
        fmt.Printf("Unsubscribed from %s with ID %d\n", eventType, id)
    }
}

func (em *EventManager) Publish(eventType string, data interface{}) {
    em.mu.RLock()
    defer em.mu.RUnlock()

    if callbacks, exists := em.subscribers[eventType]; exists {
        for _, callback := range callbacks {
            callback(data)
        }
    }
}