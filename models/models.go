package models

import (
	"fmt"
	"sync"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ReferenceValues struct {
	Models []string `json:"models"`
	Techs  []string `json:"techs"`
	Status []string `json:"status"`
}

type Robot struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Model       string   `json:"model"`
	Tech        []string `json:"tech"`
	Status      string   `json:"status"`
}

type RobotStore struct {
	mu     sync.RWMutex
	Robots map[string]Robot
	Refs   ReferenceValues
}

func NewRobotStore() *RobotStore {
	return &RobotStore{
		Robots: make(map[string]Robot),
		Refs: ReferenceValues{
			Models: []string{"car", "humanoid", "transformation"},
			Techs:  []string{"AI", "car", "robot", "cyborg", "humanoid"},
			Status: []string{"progress", "active", "inactive"},
		},
	}
}

func (rs *RobotStore) LockRead() {
	rs.mu.RLock()
}

func (rs *RobotStore) UnlockRead() {
	rs.mu.RUnlock()
}

func (rs *RobotStore) Lock() {
	rs.mu.Lock()
}

// Unlock write
func (rs *RobotStore) Unlock() {
	rs.mu.Unlock()
}

func (rs *RobotStore) GetRobot(code string) (Robot, bool) {
	robot, exists := rs.Robots[code]
	return robot, exists
}

func (rs *RobotStore) ExistsRobot(code string) bool {
	_, exists := rs.Robots[code]
	return exists
}

func (rs *RobotStore) AddRobot(robot Robot) error {
	// validation code cannot be empty
	if robot.Code == "" {
		return &ValidationError{
			Field:   "Code",
			Message: "Code cannot be empty",
		}
	}

	// Validasi model
	if !contains(rs.Refs.Models, robot.Model) {
		return &ValidationError{
			Field:   "Model",
			Message: fmt.Sprintf("Invalid model. Allowed models: %v", rs.Refs.Models),
		}
	}

	// Validasi tech
	for _, tech := range robot.Tech {
		if !contains(rs.Refs.Techs, tech) {
			return &ValidationError{
				Field:   "Tech",
				Message: fmt.Sprintf("Invalid tech: %s. Allowed techs: %v", tech, rs.Refs.Techs),
			}
		}
	}

	if !contains(rs.Refs.Status, robot.Status) {
		return &ValidationError{
			Field:   "Status",
			Message: fmt.Sprintf("Invalid status. Allowed status: %v", rs.Refs.Status),
		}
	}

	if robot.Name == "" {
		return &ValidationError{
			Field:   "Name",
			Message: "Name cannot be empty",
		}
	}

	rs.Robots[robot.Code] = robot
	return nil
}

func (rs *RobotStore) DeleteRobot(code string) {
	delete(rs.Robots, code)
}

func (rs *RobotStore) UpdateReferenceValues(newRefs ReferenceValues) {
	for _, model := range newRefs.Models {
		if !contains(rs.Refs.Models, model) {
			rs.Refs.Models = append(rs.Refs.Models, model)
		}
	}

	for _, tech := range newRefs.Techs {
		if !contains(rs.Refs.Techs, tech) {
			rs.Refs.Techs = append(rs.Refs.Techs, tech)
		}
	}

	for _, status := range newRefs.Status {
		if !contains(rs.Refs.Status, status) {
			rs.Refs.Status = append(rs.Refs.Status, status)
		}
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
