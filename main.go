package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	//"log"
	"log/slog"
	"net/http"
	"strconv"

	//"github.com/gorilla/mux"
	"github.com/sahasajib/exp-1/types"
)


var students[] types.Student

// var db *sql.DB

// func initDB(){
// 	var err error
// 	connStr := "host=localhost port=5432 user=postgre password=12345678 dbname=mydb sslmode=disable"

// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil{
// 		log.Fatal("Failed to connect to database", "error", err)
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Failed to ping database", "error", err)
// 	}
// 	slog.Info("Connected to database successfully")
// 	// Initialize students slice with some data
// }

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
	student.ID = int64(len(students) + 1 )
	students = append(students, student)
	sendData(w, student, http.StatusCreated)
}
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handleOptions(w, r)
	if r.Method != http.MethodDelete{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
		// Get ID from path param
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}	
	for i, s := range students {
		fmt.Fprintln(w, "checking student ID: ",s.ID)
		if s.ID == id {
			students = append(students[:i], students[i+1:]...)
			slog.Info("Student deleted", "id", id)
			sendData(w, s, http.StatusOK)
			return
		} 
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handleOptions(w, r)
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var studentUp types.Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&studentUp)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, s := range students {
		if s.ID == studentUp.ID {
			students[i] = studentUp
			slog.Info("Student updated", "id", studentUp.ID)
			sendData(w, studentUp, http.StatusOK)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/student", getStudent)
	router.HandleFunc("/student/create", createStudent)
	router.HandleFunc("/student/delete", deleteStudent)
	router.HandleFunc("/student/update", updateStudent)

	slog.Info("Server is running on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}

