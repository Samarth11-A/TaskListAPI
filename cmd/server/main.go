package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Samarth11-A/TaskListAPI/internal/config"
	"github.com/Samarth11-A/TaskListAPI/internal/database"
	"github.com/Samarth11-A/TaskListAPI/internal/models"
	pb "github.com/Samarth11-A/TaskList_proto/api"
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

// CreateTask creates a new task and adds it to the database
func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	log.Printf("Received CreateTask request: %v", req)

	// Convert protobuf request to internal model
	createReq := models.FromProtoCreateTaskRequest(req)

	// Validate the request
	if err := createReq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Create internal task model
	now := time.Now()
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       createReq.Title,
		Description: createReq.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Store the task
	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		log.Printf("Failed to create task: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	log.Printf("Created task with ID: %s", task.ID)
	return task.ToProtoCreateTaskResponse(), nil
}

// GetTask retrieves a task by ID
func (s *server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	log.Printf("Received GetTask request: %v", req)

	task, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	return task.ToProtoGetTaskResponse(), nil
}

// ListTasks returns a list of all tasks
func (s *server) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	log.Printf("Received ListTasks request: %v", req)

	// Convert protobuf request to internal model
	listReq := models.FromProtoListTasksRequest(req)

	tasks, err := s.taskRepo.ListTasks(ctx, listReq)
	if err != nil {
		log.Printf("Failed to list tasks: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	return tasks.ToProtoListTasksResponse(), nil
}

// UpdateTask updates an existing task
func (s *server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	log.Printf("Received UpdateTask request: %v", req)

	// Convert protobuf request to internal model
	updateReq := models.FromProtoUpdateTaskRequest(req)

	// Validate the request
	if err := updateReq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Check if task exists and get current task
	existingTask, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	// Update task fields
	existingTask.Title = updateReq.Title
	existingTask.Description = updateReq.Description
	existingTask.Completed = updateReq.Completed
	existingTask.UpdatedAt = time.Now()

	// Store updated task
	if err := s.taskRepo.UpdateTask(ctx, existingTask); err != nil {
		log.Printf("Failed to update task: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return existingTask.ToProtoUpdateTaskResponse(), nil
}

// DeleteTask removes a task by ID
func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	log.Printf("Received DeleteTask request: %v", req)

	// Delete task
	if err := s.taskRepo.DeleteTask(ctx, req.Id); err != nil {
		log.Printf("Failed to delete task: %v", err)
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	return models.ToProtoDeleteTaskResponse(true), nil
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
