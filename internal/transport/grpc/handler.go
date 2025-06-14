package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	taskpb "github.com/4nar1k/project-protos/proto/task"
	userpb "github.com/4nar1k/project-protos/proto/user"
	tasksvc "github.com/4nar1k/tasks-service/internal/task"
)

type Handler struct {
	svc        *tasksvc.TaskService
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *tasksvc.TaskService, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	// Проверяем, что title и user_id не пустые
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}
	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	// Проверяем пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.GetUserId()}); err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", req.GetUserId(), err)
	}

	// Внутренняя логика
	t, err := h.svc.CreateTask(tasksvc.Task{UserID: req.GetUserId(), Task: req.GetTitle()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	// Ответ
	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			IsDone: t.IsDone,
			UserId: t.UserID,
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	tasks, err := h.svc.GetTasksByUserID(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tasks: %v", err)
	}

	for _, t := range tasks {
		if uint32(t.ID) == req.GetId() {
			return &taskpb.Task{
				Id:     uint32(t.ID),
				Title:  t.Task,
				IsDone: t.IsDone,
				UserId: t.UserID,
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "task with id %d not found", req.GetId())
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	// Проверяем, что user_id указан
	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	// Получаем задачи пользователя без проверки существования пользователя
	tasks, err := h.svc.GetTasksByUserID(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	// Преобразуем в proto
	protoTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		protoTasks = append(protoTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			IsDone: t.IsDone,
			UserId: t.UserID,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: protoTasks}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	// Проверяем, что id и title указаны
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}

	// Получаем задачу по ID, чтобы проверить user_id
	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tasks: %v", err)
	}
	var userID uint32
	taskExists := false
	for _, t := range tasks {
		if uint32(t.ID) == req.GetId() {
			userID = t.UserID
			taskExists = true
			break
		}
	}
	if !taskExists {
		return nil, status.Errorf(codes.NotFound, "task with id %d not found", req.GetId())
	}

	// Проверяем пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: userID}); err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", userID, err)
	}

	// Обновляем задачу
	t, err := h.svc.UpdateTaskByID(req.GetId(), tasksvc.Task{
		Task:   req.GetTitle(),
		IsDone: req.GetIsDone(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	// Ответ
	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			IsDone: t.IsDone,
			UserId: t.UserID,
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	// Проверяем, что id указан
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// Получаем задачу по ID, чтобы проверить user_id
	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tasks: %v", err)
	}
	var userID uint32
	taskExists := false
	for _, t := range tasks {
		if uint32(t.ID) == req.GetId() {
			userID = t.UserID
			taskExists = true
			break
		}
	}
	if !taskExists {
		return nil, status.Errorf(codes.NotFound, "task with id %d not found", req.GetId())
	}

	// Проверяем пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: userID}); err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", userID, err)
	}

	// Удаляем задачу
	err = h.svc.DeleteTaskByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	// Ответ
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}
