package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/Samarth11-A/TaskListAPI/api/proto"
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
	// Setup a db

	tasks map[string]*pb.Task
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
	s.tasks[task.Id] = task

	log.Printf("Created task with ID: %s", task.Id)
	return task, nil
}

// GetTask retrieves a task by ID
func (s *server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	log.Printf("Received GetTask request: %v", req)

	// Check if task exists
	task, exists := s.tasks[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	return task, nil
}

// ListTasks returns a list of all tasks
func (s *server) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	log.Printf("Received ListTasks request: %v", req)

	// Simple implementation without pagination for now
	var tasks []*pb.Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return &pb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

// UpdateTask updates an existing task
func (s *server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	log.Printf("Received UpdateTask request: %v", req)

	// Check if task exists
	task, exists := s.tasks[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	// Update task fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	task.Completed = req.Completed
	task.UpdatedAt = time.Now().Format(time.RFC3339)

	// Store updated task
	s.tasks[req.Id] = task

	return task, nil
}

// DeleteTask removes a task by ID
func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	log.Printf("Received DeleteTask request: %v", req)

	// Check if task exists
	_, exists := s.tasks[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "task with ID %s not found", req.Id)
	}

	// Delete task
	delete(s.tasks, req.Id)

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	taskServer := &server{
		tasks: make(map[string]*pb.Task),
	}
	pb.RegisterTaskListServer(s, taskServer)

	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
