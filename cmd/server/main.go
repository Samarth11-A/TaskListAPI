package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/Samarth11-A/TaskListAPI/api/proto"
	"github.com/Samarth11-A/TaskListAPI/internal/config"
	"github.com/Samarth11-A/TaskListAPI/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

const (
	port = ":50051"
)

// server is used to implement the TaskList service
type server struct {
	pb.UnimplementedTaskListServer
	taskRepo *database.TaskRepository
}

// CreateTask creates a new task and adds it to the in-memory store
func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	log.Printf("Received CreateTask request: %v", req)

	// Basic validation
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title cannot be empty")
	}

	// Create a new task
	now := time.Now().Format(time.RFC3339)
	task := &pb.Task{
		Id:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Store the task
	// s.tasks[task.Id] = task

	// log.Printf("Created task with ID: %s", task.Id)
	// return task, nil
	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		log.Printf("Failed to create task: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	log.Printf("Created task with ID: %s", task.Id)
	return task, nil
}

// GetTask retrieves a task by ID
func (s *server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	log.Printf("Received GetTask request: %v", req)

	task, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	return task, nil
}

// ListTasks returns a list of all tasks
func (s *server) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	log.Printf("Received ListTasks request: %v", req)

	tasks, err := s.taskRepo.ListTasks(ctx)
	if err != nil {
		log.Printf("Failed to list tasks: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	return &pb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

// UpdateTask updates an existing task
func (s *server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	log.Printf("Received UpdateTask request: %v", req)

	// Check if task exists
	existingTask, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	// Update task fields
	if req.Title != "" {
		existingTask.Title = req.Title
	}
	if req.Description != "" {
		existingTask.Description = req.Description
	}
	existingTask.Completed = req.Completed
	existingTask.UpdatedAt = time.Now().Format(time.RFC3339)

	// Store updated task
	if err := s.taskRepo.UpdateTask(ctx, existingTask); err != nil {
		log.Printf("Failed to update task: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return existingTask, nil
}

// DeleteTask removes a task by ID
func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	log.Printf("Received DeleteTask request: %v", req)

	// Delete task
	if err := s.taskRepo.DeleteTask(ctx, req.Id); err != nil {
		log.Printf("Failed to delete task: %v", err)
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL connection
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tasks table if it doesn't exist
	if err := db.CreateTasksTable(); err != nil {
		log.Fatalf("Failed to create tasks table: %v", err)
	}

	// Create task repository
	taskRepo := database.NewTaskRepository(db)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	taskServer := &server{
		taskRepo: taskRepo,
	}
	pb.RegisterTaskListServer(s, taskServer)

	log.Printf("Server listening on port %s", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
