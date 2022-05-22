package services

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Task struct {
	Id          uuid.UUID  `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Assigned    uuid.UUID  `json:"assigned"`
}

type TaskStatus int

const (
	TODO TaskStatus = iota
	DONE
)

type TaskService interface {
	CreateTask(description string, assignee uuid.UUID) Task
	GetTask(id uuid.UUID) (Task, error)
	GetAllUserTasks(assignee uuid.UUID) []Task
	GetAllTasks() []Task
	CompleteTask(taskId uuid.UUID, userId uuid.UUID) (Task, error)
}

type taskService struct {
	sync.Mutex
	tasks map[uuid.UUID]Task
}

func NewTaskService() TaskService {
	return &taskService{tasks: make(map[uuid.UUID]Task)}
}

func (ts *taskService) CreateTask(description string, assignee uuid.UUID) Task {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:          uuid.New(),
		Description: description,
		Assigned:    assignee,
	}

	ts.tasks[task.Id] = task
	return task
}

func (ts *taskService) GetTask(id uuid.UUID) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%s not found", id.String())
	}
}

func (ts *taskService) GetAllUserTasks(assignee uuid.UUID) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		if task.Assigned == assignee {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (ts *taskService) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (ts *taskService) CompleteTask(taskId uuid.UUID, userId uuid.UUID) (Task, error) {
	var task, err = ts.GetTask(taskId)
	if err != nil {
		return task, err
	}

	if task.Assigned != userId {
		return task, fmt.Errorf("task with id=%d does not assigned to user=%d", taskId, userId)
	}

	if task.Status == DONE {
		return task, fmt.Errorf("task with taskId=%d is already completed", taskId)
	}

	ts.Lock()
	defer ts.Unlock()

	task.Status = DONE
	ts.tasks[task.Id] = task
	return task, nil
}
