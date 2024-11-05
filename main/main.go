package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/manager"
	"github.com/kiquetal/orchestration-go-scratch/node"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"github.com/kiquetal/orchestration-go-scratch/worker"
	"os"
	"time"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task 1",
		State:  task.Pending,
		Image:  "nginx",
		Memory: 512,
		Disk:   512,
	}
	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Task:      t,
		Timestamp: time.Now(),
	}
	fmt.Printf("Task: %+v\n", t)
	fmt.Printf("Task Event: %+v\n", te)

	w := worker.Worker{
		Name:  "Worker 1",
		Queue: queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("Worker: %+v\n", w)
	w.CollectStats()
	w.RunTask(&t)
	w.StopTask(&t)

	m := manager.Manager{
		Pending: queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}
	fmt.Printf("Manager: %+v\n", m)
	m.SelectWorker()
	m.UpdateTask()
	m.SendWork()

	n := node.Node{
		Name:   "Node 1",
		Ip:     "192,168.1.1",
		Cores:  2,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}

	fmt.Printf("Node: %+v\n", n)
	fmt.Printf("create a container\n")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("Error: %v\n", createResult.Error)
		os.Exit(1)
	}
	fmt.Printf("stop a container\n")
	time.Sleep(10 * time.Second)
	stopContainer(dockerTask, createResult.ContainerId)

}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user"},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()
	if result.Error != nil {
		return &d, &result
	}
	fmt.Printf("Container %s is running with %v values \n", result.Container, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		return &result
	}
	fmt.Printf("Container %s is stopped \n", result.Container)
	return &result
}
