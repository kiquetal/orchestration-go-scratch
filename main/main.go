package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"github.com/kiquetal/orchestration-go-scratch/worker"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	host := os.Getenv("ORCHESTRATION_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ORCHESTRATION_PORT")
	if port == "" {
		port = "8080"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to convert port to integer: %v", err)
	}
	fmt.Println("Starting the orchestration system")
	w := worker.Worker{
		Name:  "worker-1",
		Queue: queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{
		Worker:  &w,
		Address: host,
		Port:    portInt,
	}
	go runTasks(&w)
	go w.CollectStats()
	api.Start()

}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				fmt.Println("Error running task: ", result.Error)
			}

		} else {
			fmt.Println("No tasks in queue")
		}
		log.Printf("Worker %s is running %d tasks", w.Name, w.Queue.Len())
		time.Sleep(10 * time.Second)
	}
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
