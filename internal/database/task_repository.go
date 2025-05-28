package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "github.com/Samarth11-A/TaskListAPI/api/proto"
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
func (r *TaskRepository) CreateTask(ctx context.Context, task *pb.Task) error {
	query := `
    INSERT INTO tasks (id, title, description, completed, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		task.Id, task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetTask retrieves a task by ID
func (r *TaskRepository) GetTask(ctx context.Context, id string) (*pb.Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = $1`

	var task pb.Task
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// ListTasks retrieves all tasks
func (r *TaskRepository) ListTasks(ctx context.Context) ([]*pb.Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks`

	var tasks []*pb.Task
	err := r.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (r *TaskRepository) UpdateTask(ctx context.Context, task *pb.Task) error {
	query := `
    UPDATE tasks 
    SET title = $2, description = $3, completed = $4, updated_at = $5
    WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		task.Id, task.Title, task.Description, task.Completed, task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found with ID: %s", task.Id)
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
