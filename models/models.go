package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"passowrd"`
}

type Event struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type Task struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	IsDone   bool   `json:"is_done"`
	EventId  uint   `json:"event_id"`
	TaskType string `json:"task_type"`
}
