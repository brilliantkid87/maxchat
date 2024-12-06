package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"robot/models"
	"strings"
)

var robotStore *models.RobotStore

func main() {
	robotStore = models.NewRobotStore()
	loadInitialData()

	http.HandleFunc("/robots", handleRobots)
	http.HandleFunc("/robots/", handleRobotByCode)
	http.HandleFunc("/references", handleReferences)
	http.HandleFunc("/references/update", handleReferencesUpdate)

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadInitialData() {
	file, err := os.Open("data/initial_data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data struct {
		Robots []models.Robot `json:"robots"`
	}

	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &data)

	for _, robot := range data.Robots {
		robotStore.AddRobot(robot)
	}
}

func handleRobots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRobots(w, r)
	case http.MethodPost:
		createRobot(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRobotByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/robots/")
	switch r.Method {
	case http.MethodGet:
		getRobotByCode(w, r, code)
	case http.MethodPut:
		updateRobot(w, r, code)
	case http.MethodDelete:
		deleteRobot(w, r, code)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRobots(w http.ResponseWriter, r *http.Request) {
	robotStore.LockRead()
	defer robotStore.UnlockRead()

	model := r.URL.Query().Get("model")
	techQuery := r.URL.Query().Get("tech")
	techs := strings.Split(techQuery, ",")

	var filteredRobots []models.Robot

	for _, robot := range robotStore.Robots {
		// Filter by model
		if model != "" && robot.Model != model {
			continue
		}

		// Multi-tech filter
		if techQuery != "" {
			match := true
			for _, tech := range techs {
				if !containsTech(robot.Tech, tech) {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		filteredRobots = append(filteredRobots, robot)
	}

	json.NewEncoder(w).Encode(filteredRobots)
}

func containsTech(techs []string, searchTech string) bool {
	for _, tech := range techs {
		if tech == searchTech {
			return true
		}
	}
	return false
}

func handleReferences(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(robotStore.Refs)
}

func handleReferencesUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	robotStore.Lock()
	defer robotStore.Unlock()

	var newRefs models.ReferenceValues
	err := json.NewDecoder(r.Body).Decode(&newRefs)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	robotStore.UpdateReferenceValues(newRefs)
	json.NewEncoder(w).Encode(robotStore.Refs)
}

func createRobot(w http.ResponseWriter, r *http.Request) {
	var robot models.Robot
	err := json.NewDecoder(r.Body).Decode(&robot)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	robotStore.Lock()
	defer robotStore.Unlock()

	if robotStore.ExistsRobot(robot.Code) {
		http.Error(w, "Robot code already exists", http.StatusConflict)
		return
	}

	err = robotStore.AddRobot(robot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(robot)
}

func getRobotByCode(w http.ResponseWriter, r *http.Request, code string) {
	robotStore.LockRead()
	defer robotStore.UnlockRead()

	robot, exists := robotStore.GetRobot(code)
	if !exists {
		http.Error(w, "Robot not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(robot)
}

func updateRobot(w http.ResponseWriter, r *http.Request, code string) {
	robotStore.Lock()
	defer robotStore.Unlock()

	if !robotStore.ExistsRobot(code) {
		http.Error(w, "Robot not found", http.StatusNotFound)
		return
	}

	var updatedRobot models.Robot
	err := json.NewDecoder(r.Body).Decode(&updatedRobot)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedRobot.Code = code

	err = robotStore.AddRobot(updatedRobot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updatedRobot)
}

func deleteRobot(w http.ResponseWriter, r *http.Request, code string) {
	robotStore.Lock()
	defer robotStore.Unlock()

	if !robotStore.ExistsRobot(code) {
		http.Error(w, "Robot not found", http.StatusNotFound)
		return
	}

	// remove robot
	robotStore.DeleteRobot(code)
	w.WriteHeader(http.StatusNoContent)
}
