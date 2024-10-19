// File: internal/ecs/entity.go
package ecs

type Entity uint32

const (
    IndexBits   = 24
    VersionBits = 8
    IndexMask   = (1 << IndexBits) - 1
    VersionMask = (1 << VersionBits) - 1
)