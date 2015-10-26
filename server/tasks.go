package api

import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
)

const TasksPrefix = "/tasks/"

func addApiTasks(r *mux.Router) {
	r.HandleFunc(TasksPrefix, wrapJsonHandler(ListTasks)).Methods("GET")
	r.HandleFunc(TasksPrefix, wrapJsonHandler(NewTask)).Methods("POST", "OPTIONS")
	r.HandleFunc(TasksPrefix+"{id}", wrapJsonHandler(GetTask)).Methods("GET")
	r.HandleFunc(TasksPrefix+"{id}", wrapJsonHandler(UpdateTask)).Methods("PUT")
	r.HandleFunc(TasksPrefix+"{id}", wrapJsonHandler(DeleteTask)).Methods("DELETE", "OPTIONS")
}

type TaskPriority string

const (
	prio_normal  TaskPriority = "normal"
	prio_warning TaskPriority = "warning"
	prio_asap    TaskPriority = "asap"
)

type Task struct {
	ID       int64
	Title    string
	Done     bool
	Priority TaskPriority
}

func newDefaultTask() *Task {
	return &Task{0, "", false, prio_normal}
}

// req: GET /tasks/
// resp: 200
func ListTasks(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	var tasks []Task
	if _, err := datastore.NewQuery("Task").GetAll(c, &tasks); err != nil {
		return err
	}
	if tasks == nil {
		tasks = []Task{}
	}

	return json.NewEncoder(w).Encode(tasks)
}

// req: POST /tasks/ {"Title": "Buy bread"}
// resp: 201
func NewTask(w http.ResponseWriter, r *http.Request) error {
	req := struct{ Title string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	if req.Title == "" {
		return fmt.Errorf("Empty title!")
	}

	newTask := newDefaultTask()
	newTask.Title = req.Title

	c := appengine.NewContext(r)
	taskKey, e := datastore.Put(c, datastore.NewIncompleteKey(c, "Task", nil), newTask)
	if e != nil {
		return e
	}

	newTask.ID = taskKey.IntID()
	_, e = datastore.Put(c, taskKey, newTask)
	// log.Println("newTask.ID -> ", newTask.ID)
	if e != nil {
		return e
	}

	newUrl := r.URL.Path + strconv.FormatInt(newTask.ID, 10)
	w.Header().Set("Location", newUrl)
	w.WriteHeader(http.StatusCreated)
	return nil
}

// req: /tasks/99
// resp: 200
func GetTask(w http.ResponseWriter, r *http.Request) error {
	taskId, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}

	c := appengine.NewContext(r)
	taskKey := datastore.NewKey(c, "Task", "", taskId, nil)

	var task Task
	err = datastore.Get(c, taskKey, &task)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(task)
}

// req: /tasks/99 {ID:99, Title:"title",...}
// resp: 204
func UpdateTask(w http.ResponseWriter, r *http.Request) error {
	taskId, err := parseID(r)
	// log.Println("taskId -> ", taskId)
	if err != nil {
		return badRequest{err}
	}
	var taskUI Task
	if err := json.NewDecoder(r.Body).Decode(&taskUI); err != nil {
		return badRequest{err}
	}
	if taskUI.ID != taskId {
		return badRequest{fmt.Errorf("Inconsistent task IDs")}
	}

	c := appengine.NewContext(r)
	taskKey := datastore.NewKey(c, "Task", "", taskId, nil)

	newTask := &Task{taskId, taskUI.Title, taskUI.Done, taskUI.Priority}
	_, err = datastore.Put(c, taskKey, newTask)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// req: /tasks/99
// resp: 204
func DeleteTask(w http.ResponseWriter, r *http.Request) error {
	taskId, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}

	c := appengine.NewContext(r)
	taskKey := datastore.NewKey(c, "Task", "", taskId, nil)

	err = datastore.Delete(c, taskKey)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
