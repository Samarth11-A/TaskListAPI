package models

import (
	"time"

	pb "github.com/Samarth11-A/TaskList_proto/api"
)

// ToProtoTask converts an internal Task to a protobuf Task
func (t *Task) ToProtoTask() *pb.Task {
	return &pb.Task{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

// FromProtoTask converts a protobuf Task to an internal Task
func FromProtoTask(protoTask *pb.Task) (*Task, error) {
	createdAt, err := time.Parse(time.RFC3339, protoTask.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, protoTask.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:          protoTask.Id,
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Completed:   protoTask.Completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// FromProtoCreateTaskRequest converts a protobuf CreateTaskRequest to internal type
func FromProtoCreateTaskRequest(req *pb.CreateTaskRequest) *CreateTaskRequest {
	return &CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}
}

// FromProtoUpdateTaskRequest converts a protobuf UpdateTaskRequest to internal type
func FromProtoUpdateTaskRequest(req *pb.UpdateTaskRequest) *UpdateTaskRequest {
	return &UpdateTaskRequest{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}
}

// FromProtoListTasksRequest converts a protobuf ListTasksRequest to internal type
func FromProtoListTasksRequest(req *pb.ListTasksRequest) *ListTasksRequest {
	return &ListTasksRequest{
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}
}

// ToProtoCreateTaskResponse converts internal data to protobuf CreateTaskResponse
func (t *Task) ToProtoCreateTaskResponse() *pb.CreateTaskResponse {
	return &pb.CreateTaskResponse{
		Task: t.ToProtoTask(),
	}
}

// ToProtoGetTaskResponse converts internal Task to protobuf GetTaskResponse
func (t *Task) ToProtoGetTaskResponse() *pb.GetTaskResponse {
	return &pb.GetTaskResponse{
		Task: t.ToProtoTask(),
	}
}

// ToProtoUpdateTaskResponse converts internal Task to protobuf UpdateTaskResponse
func (t *Task) ToProtoUpdateTaskResponse() *pb.UpdateTaskResponse {
	return &pb.UpdateTaskResponse{
		Task: t.ToProtoTask(),
	}
}

// ToProtoListTasksResponse converts internal ListTasksResponse to protobuf
func (r *ListTasksResponse) ToProtoListTasksResponse() *pb.ListTasksResponse {
	protoTasks := make([]*pb.Task, len(r.Tasks))
	for i, task := range r.Tasks {
		protoTasks[i] = task.ToProtoTask()
	}

	return &pb.ListTasksResponse{
		Tasks:         protoTasks,
		NextPageToken: r.NextPageToken,
	}
}

// ToProtoDeleteTaskResponse creates a protobuf DeleteTaskResponse
func ToProtoDeleteTaskResponse(success bool) *pb.DeleteTaskResponse {
	return &pb.DeleteTaskResponse{
		Success: success,
	}
}