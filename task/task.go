package task

import (
	"context"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"io"
	"os"
	"time"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Cpu           float64
	Image         string
	Memory        int64
	Disk          int64
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

type Docker struct {
	Client *client.Client
	Config Config
}
type Config struct {
	// Name of the task, also used as the container name
	Name string
	// AttachStdin boolean which determines if stdin should be attached
	AttachStdin bool
	// AttachStdout boolean which determines if stdout should be attached
	AttachStdout bool
	// AttachStderr boolean which determines if stderr should be attached
	AttachStderr bool
	// ExposedPorts list of ports exposed
	ExposedPorts nat.PortSet
	// Cmd to be run inside container (optional)
	Cmd []string
	// Image used to run the container
	Image string
	// Cpu
	Cpu float64
	// Memory in MiB
	Memory int64
	// Disk in GiB
	Disk int64
	// Env variables
	Env []string
	// RestartPolicy for the container ["", "always", "unless-stopped", "on-failure"]
	RestartPolicy string
}

type DockerResult struct {
	Error     error
	Action    string
	Container string
	Result    string
}

func (t *Task) NewConfig() *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Cpu:           t.Cpu,
		Memory:        t.Memory,
		Disk:          t.Disk,
		RestartPolicy: t.RestartPolicy,
	}
}

func NewDocker(c *Config) *Docker {
	client_docker, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{
		Config: *c,
		Client: client_docker,
	}
}
func (d *Docker) Run() DockerResult {

	ctx := context.Background()
	reader, err := d.Client.ImagePull(
		ctx,
		d.Config.Image,
		image.PullOptions{})

	if err != nil {
		return DockerResult{
			Error:     err,
			Action:    "Pull",
			Container: d.Config.Name,
			Result:    "Failed",
		}
	}
	io.Copy(os.Stdout, reader)
	defer reader.Close()
	return DockerResult{
		Error:  nil,
		Action: "Pull",
	}

}
