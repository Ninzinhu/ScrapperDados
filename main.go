package main

import (
    "encoding/json"
    "net/http"
    "sync"
)

type Task struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var tasks = make(map[string]Task)
var mu sync.Mutex

// Criar nova Task
func createTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil || task.ID == "" || task.Name == "" {
        http.Error(w, "Input Inválido", http.StatusBadRequest)
        return
    }

    mu.Lock()
    tasks[task.ID] = task
    mu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

// Ler todas Tasks
func readTasks(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    taskList := make([]Task, 0, len(tasks))
    for _, task := range tasks {
        taskList = append(taskList, task)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(taskList)
}

// Atualizar a Task
func updateTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil || task.ID == "" || task.Name == "" {
        http.Error(w, "Input Inválido", http.StatusBadRequest)
        return
    }

    mu.Lock()
    if _, exists := tasks[task.ID]; !exists {
        mu.Unlock()
        http.Error(w, "Task não encontrada", http.StatusNotFound)
        return
    }
    tasks[task.ID] = task
    mu.Unlock()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(task)
}

// Deletar a task
func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    mu.Lock()
    defer mu.Unlock()

    if _, exists := tasks[id]; !exists {
        http.Error(w, "Task não encontrada", http.StatusNotFound)
        return
    }
    delete(tasks, id)

    w.WriteHeader(http.StatusNoContent)
}

func main() {
    http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            createTask(w, r)
        case http.MethodGet:
            readTasks(w, r)
        case http.MethodPut:
            updateTask(w, r)
        case http.MethodDelete:
            deleteTask(w, r)
        default:
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        }
    })

    http.ListenAndServe(":8080", nil)
}