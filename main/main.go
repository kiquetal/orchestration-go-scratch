package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"github.com/kiquetal/orchestration-go-scratch/worker"
	"time"
)

func main() {

	m := make(map[uuid.UUID]*task.Task)
	taskE := &task.Task{
		ID:    uuid.New(),
		Name:  "test-task-1",
		State: task.Scheduled,
		Cpu:   0.5,
		Image: "strm/helloworld-http",
	}
	m[taskE.ID] = taskE

	w := worker.Worker{
		Name:  "worker-1",
		Queue: queue.New(),
		Db:    m,
	}
	w.AddTask(*taskE)
	result := w.RunTask()
	fmt.Println(result)
	if result.Error != nil {
		panic(result.Error)
	}
	taskE.ContainerID = result.Container
	fmt.Println("Task is running with container id: ", taskE.ContainerID)
	fmt.Println("Sleep 7 seconds")
	time.Sleep(7 * time.Second)
	taskE.State = task.Completed
	w.AddTask(*taskE)
	result = w.RunTask()
	fmt.Println(result)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("Task is stopped")
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
