package main

import (
	//	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/sirupsen/logrus"
	// "github.com/spf13/viper"
)

// Response struct
type Response struct {
	Message string `json:"message"`
}

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var tasks = []Task{}
var idCounter = 1

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	task.ID = idCounter
	idCounter++
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			var updatedTask Task
			err := json.NewDecoder(r.Body).Decode(&updatedTask)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			updatedTask.ID = task.ID
			tasks[i] = updatedTask
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Prometheus metrics
var httpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests",
}, []string{"method", "path"},
)

// func setupConfig() {
// 	viper.SetConfigName("config")
// 	viper.AddConfigPath(".")

// 	viper.SetDefault("server.port", "8080")

// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Printf("Config file not found, using defaults: %s", err)
// 	}
// }

// func setupLogger() *logrus.Logger {
// 	log := logrus.New()
// 	log.SetFormatter(&logrus.JSONFormatter{})
// 	return log
// }

func handler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
	response := Response{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func setupDatabase() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "./data.db")
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, name TEXT)")
// 	return db, err
// }

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {

	// Setup Configuration and Logger

	// setupConfig()
	// logger := setupLogger()

	// Register metrics

	prometheus.MustRegister(httpRequests)

	// API handler
	// Endpoints
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTasks(w, r)
		} else if r.Method == "POST" {
			createTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			updateTask(w, r)
		} else if r.Method == "DELETE" {
			deleteTask(w, r)
		} else {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", handler)

	// http.HandleFunc("/tasks", createTask)      // POST
	// http.HandleFunc("/tasks", getTasks)        // GET
	// http.HandleFunc("/tasks/{id}", updateTask) // PUT
	// http.HandleFunc("/tasks/{id}", deleteTask) //DELETE

	// http.Handle("/", loggingMiddleware(http.HandlerFunc(handler)))

	// Prometheus Metrics Endpoint
	http.Handle("/metrics", promhttp.Handler())

	// port := viper.GetString("server.port")
	// logger.Infof("Starting server on port %s", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	logger.Fatalf("could not start server: %v", err)
	// }

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}
