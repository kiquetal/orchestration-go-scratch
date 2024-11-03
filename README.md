#### Trying to understand the orchestration nuts & bolts





#### Docker commands

- Start a container: `docker run -d -p 80:80 --name webserver nginx`
- Stop a container: `docker stop webserver`
- Inspect a container: `docker inspect webserver`
- Remove a container: `docker rm webserver`
- Remove an image: `docker rmi nginx`
- List containers: `docker ps -a`
- List images: `docker images`
- List networks: `docker network ls`
- List volumes: `docker volume ls`
- List all containers: `docker ps -a`
- List all containers with their IDs: `docker ps -aq`
- Remove specific containers: `docker rm container_id`

```mermaid

sequenceDiagram
    participant U as User
    participant M as Manager
    participant S as Scheduler
    participant W as Worker
    participant D as Docker Container

    U->>M: Send request to start/stop task
    M->>S: Consult for available worker(s)
    S-->>M: Return available worker(s)
    M->>W: Assign task
    W->>D: Run task
    D-->>W: Update task state
    W->>M: Report statistics
    M->>M: Update records of tasks and workers

```
