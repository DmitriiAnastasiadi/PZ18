package service

import (
	"strconv"

	"sync"
)

type Task struct {
	ID string `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	DueDate string `json:"due_date"`

	Done bool `json:"done"`
}

var tasks = make(map[string]Task)

var mu sync.Mutex

var counter int

func CreateTask(title, desc, due string) Task {

	mu.Lock()

	defer mu.Unlock()

	counter++

	id := "t_" + strconv.Itoa(counter)

	task := Task{ID: id, Title: title, Description: desc, DueDate: due, Done: false}

	tasks[id] = task

	return task

}

func GetTasks() []Task {

	mu.Lock()

	defer mu.Unlock()

	var list []Task

	for _, t := range tasks {

		list = append(list, Task{ID: t.ID, Title: t.Title, Done: t.Done})

	}

	return list

}

func GetTask(id string) (Task, bool) {

	mu.Lock()

	defer mu.Unlock()

	t, ok := tasks[id]

	return t, ok

}

func UpdateTask(id string, updates map[string]interface{}) (Task, bool) {

	mu.Lock()

	defer mu.Unlock()

	t, ok := tasks[id]

	if !ok {

		return Task{}, false

	}

	if title, ok := updates["title"].(string); ok {

		t.Title = title

	}

	if done, ok := updates["done"].(bool); ok {

		t.Done = done

	}

	tasks[id] = t

	return t, true

}

func DeleteTask(id string) bool {

	mu.Lock()

	defer mu.Unlock()

	_, ok := tasks[id]

	if ok {

		delete(tasks, id)

	}

	return ok

}
