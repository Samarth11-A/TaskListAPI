package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Samarth11-A/TaskList_proto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051" // Default gRPC server address
	timeout = 5 * time.Second
)

// Implement the client logic here which will interact with the server
// using the generated protobuf code. This will include creating a gRPC client,
// making requests to the server, and handling responses.

func main() {
	// Set up a connection to the server
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewTaskListClient(conn)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Example: Call ListTasks RPC
	listResp, err := client.ListTasks(ctx, &pb.ListTasksRequest{PageSize: 10})
	if err != nil {
		log.Fatalf("Could not list tasks: %v", err)
	}

	// Display the result
	fmt.Println("Tasks:", len(listResp.Tasks))
	for i, task := range listResp.Tasks {
		fmt.Printf("%d. %s - %s (Completed: %t)\n", i+1, task.Title, task.Description, task.Completed)
	}

	// Example usage of additional methods:
	/*
		// Create a new task
		createResp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
			Title:       "New Task",
			Description: "Task created from client",
		})
		if err != nil {
			log.Fatalf("Could not create task: %v", err)
		}
		fmt.Printf("Created Task: %s\n", createResp.Id)

		// Update a task
		updateResp, err := client.UpdateTask(ctx, &pb.UpdateTaskRequest{
			Id:          "task-id",
			Title:       "Updated Task",
			Description: "Task updated from client",
			Completed:   true,
		})
		if err != nil {
			log.Fatalf("Could not update task: %v", err)
		}
		fmt.Printf("Updated Task: %s\n", updateResp.Title)

		// Get a task
		getResp, err := client.GetTask(ctx, &pb.GetTaskRequest{Id: "task-id"})
		if err != nil {
			log.Fatalf("Could not get task: %v", err)
		}
		fmt.Printf("Retrieved Task: %s - %s\n", getResp.Title, getResp.Description)

		// Delete a task
		deleteResp, err := client.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: "task-id"})
		if err != nil {
			log.Fatalf("Could not delete task: %v", err)
		}
		fmt.Printf("Task deleted: %t\n", deleteResp.Success)
	*/
}
