import (
	"github.com/google/uuid"
	"github.com/docker/go-connections/nat"
)
type Task struct{
	ID uuid.UUID
	Name string
	State State
	Image string
	Memory int
	Disk int
	ExposedPorts nat.PortSet
	PortBindings map[string]string
	RestartPolicy string
}


