// pkg/components/lung.go

package components

type LungState int

const (
    Exhale LungState = iota
    Inhale
)

type Lung struct {
    State    LungState
    Capacity float32
    Volume   float32
}

func (l *Lung) IsComponentData() {}