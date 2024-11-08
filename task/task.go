package task

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"io"
	"log"
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
	ContainerID   string
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
	ContainerID   string
}

var TransitionMapState = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Running, Failed, Scheduled},
	Running:   {Completed, Failed, Running},
	Completed: {},
	Failed:    {},
}

type DockerResult struct {
	Error       error
	Action      string
	Container   string
	Result      string
	ContainerId string
}

func (t *Task) NewConfig() *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Cpu:           t.Cpu,
		Memory:        t.Memory,
		Disk:          t.Disk,
		ContainerID:   t.ContainerID,
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

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * 1e9),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}
	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(
		ctx,
		&cc,
		&hc,
		nil,
		nil, d.Config.Name)

	fmt.Printf("The id is %s\n", resp.ID)
	if err != nil {
		return DockerResult{
			Error:     err,
			Action:    "Create",
			Container: d.Config.Name,
			Result:    "Failed",
		}
	}

	err = d.Client.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{})

	if err != nil {
		return DockerResult{
			Error:     err,
			Action:    "Start",
			Container: d.Config.Name,
			Result:    "Failed",
		}
	}

	out, err := d.Client.ContainerLogs(
		ctx,
		resp.ID,
		container.LogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})

	if err != nil {
		return DockerResult{
			Error:     err,
			Action:    "Logs",
			Container: d.Config.Name,
			Result:    "Failed",
		}

	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		Container:   resp.ID,
		Action:      "Start",
		Result:      "Success",
		ContainerId: resp.ID,
	}

}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Stop container %s", id)
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %s", id, err)
		return DockerResult{
			Error:     err,
			Action:    "Stop",
			Container: id,
			Result:    "Failed",
		}
	}
	err = d.Client.ContainerRemove(ctx, id, container.RemoveOptions{
		Force:         false,
		RemoveLinks:   false,
		RemoveVolumes: true,
	})

	if err != nil {
		log.Printf("Error removing container %s: %s", id, err)
		return DockerResult{
			Error:     err,
			Action:    "Remove",
			Container: id,
			Result:    "Failed",
		}
	}
	return DockerResult{
		Container: id,
		Action:    "Stop",
		Result:    "Success",
	}

}
