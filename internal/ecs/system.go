package ecs

type System interface {
	Update(dt float32)
}