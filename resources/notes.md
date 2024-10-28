#### Worker responsibilities

- Run tasks as Docker containers
- Accept tasks to run from a manager
- Provide relevant statistics to the manager for the purpose of scheduling
tasks
- Keep track of its tasks and their state


#### Manager responsibilities

- Accepts requests from the user to start and stop tasks
- Schedules tasks to onto workers
- Keeps track of the state of all tasks and the machines they are running on


#### Scheduler responsibilities

- Determine a set of candidate workers on which task can run


