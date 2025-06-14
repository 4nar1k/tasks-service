package task

import (
	"context"
	"fmt"
	userpb "github.com/4nar1k/project-protos/proto/user"
)

type TaskService struct {
	repo       TaskRepository
	userClient userpb.UserServiceClient
}

func NewTaskService(repo TaskRepository, userClient userpb.UserServiceClient) *TaskService {
	return &TaskService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *TaskService) CreateTask(task Task) (Task, error) {
	if task.UserID == 0 {
		return Task{}, fmt.Errorf("user_id is required")
	}

	_, err := s.userClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: task.UserID})
	if err != nil {
		return Task{}, fmt.Errorf("user with id %d not found: %v", task.UserID, err)
	}

	return s.repo.CreateTask(task)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTasksByUserID(userID uint32) ([]Task, error) {
	return s.repo.GetTasksByUserID(userID)
}

func (s *TaskService) UpdateTaskByID(id uint32, task Task) (Task, error) {
	return s.repo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTaskByID(id uint32) error {
	return s.repo.DeleteTaskByID(id)
}
