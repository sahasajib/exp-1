package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/sahasajib/exp-1/types"
)


var students[] types.Student

func handleCors(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}
}

func sendData(w http.ResponseWriter, data interface{}, statusCode int){
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w) 
	encoder.Encode(data)
} 

func getStudent(w http.ResponseWriter, r *http.Request){
	handleCors(w)
	handleOptions(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(students) == 0 {
		http.Error(w, "No students found", http.StatusNotFound)
		return
	}

	sendData(w, students, http.StatusAccepted)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handleOptions(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var student types.Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&student)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	students = append(students, student)
	sendData(w, student, http.StatusCreated)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/student", getStudent)
	router.HandleFunc("/student/create", createStudent)

	slog.Info("Server is running on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}

