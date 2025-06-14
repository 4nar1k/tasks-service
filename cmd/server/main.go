package main

import (
	db "github.com/4nar1k/tasks-service/internal/database"
	tasksvc "github.com/4nar1k/tasks-service/internal/task"
	transportgrpc "github.com/4nar1k/tasks-service/internal/transport/grpc"
	"log"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	repo := tasksvc.NewTaskRepository(dbInstance)
	userClient, conn, err := transportgrpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to users: %v", err)
	}
	defer conn.Close()

	svc := tasksvc.NewTaskService(repo, userClient) // Добавлен userClient

	if err := transportgrpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("Tasks gRPC server error: %v", err)
	}
}
