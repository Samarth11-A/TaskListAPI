package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Samarth11-A/TaskListAPI/internal/models"
)

// TaskRepository provides methods to interact with tasks in the database
type TaskRepository struct {
	db *PostgresDB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *PostgresDB) *TaskRepository {
	return &TaskRepository{db: db}
}

// CreateTask adds a new task to the database
func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	query := `
    INSERT INTO tasks (id, title, description, completed, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Completed,
		task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339))

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetTask retrieves a task by ID
func (r *TaskRepository) GetTask(ctx context.Context, id string) (*models.Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = $1`

	var task struct {
		ID          string `db:"id"`
		Title       string `db:"title"`
		Description string `db:"description"`
		Completed   bool   `db:"completed"`
		CreatedAt   string `db:"created_at"`
		UpdatedAt   string `db:"updated_at"`
	}

	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, task.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return &models.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// ListTasks retrieves tasks with pagination support
func (r *TaskRepository) ListTasks(ctx context.Context, req *models.ListTasksRequest) (*models.ListTasksResponse, error) {
	// Set default page size if not specified
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// For simplicity, we'll just return all tasks for now
	// In a real implementation, you'd handle pagination properly
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks ORDER BY created_at DESC LIMIT $1`

	var dbTasks []struct {
		ID          string `db:"id"`
		Title       string `db:"title"`
		Description string `db:"description"`
		Completed   bool   `db:"completed"`
		CreatedAt   string `db:"created_at"`
		UpdatedAt   string `db:"updated_at"`
	}

	err := r.db.SelectContext(ctx, &dbTasks, query, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	tasks := make([]*models.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		createdAt, err := time.Parse(time.RFC3339, dbTask.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		updatedAt, err := time.Parse(time.RFC3339, dbTask.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		tasks[i] = &models.Task{
			ID:          dbTask.ID,
			Title:       dbTask.Title,
			Description: dbTask.Description,
			Completed:   dbTask.Completed,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
	}

	return &models.ListTasksResponse{
		Tasks:         tasks,
		NextPageToken: "", // Implement proper pagination tokens if needed
	}, nil
}

// UpdateTask updates an existing task
func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `
    UPDATE tasks 
    SET title = $2, description = $3, completed = $4, updated_at = $5
    WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Completed, task.UpdatedAt.Format(time.RFC3339))

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found with ID: %s", task.ID)
	}

	return nil
}

// DeleteTask removes a task by ID
func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found with ID: %s", id)
	}

	return nil
}
