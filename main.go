package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

var (
	dataList = []map[string]interface{}{}
	mu       sync.Mutex
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit-form", htmxFormSubmitHandler)
	http.HandleFunc("/delete-entry", deleteEntryHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func htmxFormSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		editIndexStr := r.FormValue("editIndex")
		editIndex, err := strconv.Atoi(editIndexStr)
		if err != nil {
			http.Error(w, "Invalid edit index", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		responseData := map[string]interface{}{
			"First Name":    r.FormValue("firstName"),
			"Last Name":     r.FormValue("lastName"),
			"Email":         r.FormValue("email"),
			"Mobile":        r.FormValue("mobile"),
			"Country":       r.FormValue("country"),
			"Date of Birth": r.FormValue("dob"),
		}

		if editIndex >= 0 && editIndex < len(dataList) {
			dataList[editIndex] = responseData
		} else {
			dataList = append(dataList, responseData)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dataList)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func deleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		editIndexStr := r.URL.Query().Get("index")
		editIndex, err := strconv.Atoi(editIndexStr)
		if err != nil {
			http.Error(w, "Invalid index", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if editIndex >= 0 && editIndex < len(dataList) {
			dataList = append(dataList[:editIndex], dataList[editIndex+1:]...)
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
