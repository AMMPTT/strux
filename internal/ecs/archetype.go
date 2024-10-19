// internal/ecs/archetype.go

package ecs

import (
	"reflect"
    "github.com/AMMPTT/strux/pkg/components"

)

type Component = components.ComponentData

type Archetype struct {
    componentTypes []reflect.Type
    entities       []Entity
    components     [][]Component
}