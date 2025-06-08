package models

import (
	"errors"
	"time"
)

// Task represents the internal domain model for a task
type Task struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the task fields
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(t.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	if len(t.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// CreateTaskRequest represents the internal request for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate validates the create task request
func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(r.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	if len(r.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// UpdateTaskRequest represents the internal request for updating a task
type UpdateTaskRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Validate validates the update task request
func (r *UpdateTaskRequest) Validate() error {
	if r.ID == "" {
		return errors.New("id cannot be empty")
	}
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(r.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	if len(r.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}
	return nil
}

// ListTasksRequest represents the internal request for listing tasks
type ListTasksRequest struct {
	PageToken string `json:"page_token"`
	PageSize  int32  `json:"page_size"`
}

// ListTasksResponse represents the internal response for listing tasks
type ListTasksResponse struct {
	Tasks         []*Task `json:"tasks"`
	NextPageToken string  `json:"next_page_token"`
}
