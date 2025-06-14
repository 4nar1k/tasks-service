package grpc

import (
	taskpb "github.com/4nar1k/project-protos/proto/task"
	userpb "github.com/4nar1k/project-protos/proto/user"
	"github.com/4nar1k/tasks-service/internal/task"
	"google.golang.org/grpc"
	"net"
)

func RunGRPC(svc *task.TaskService, uc userpb.UserServiceClient) error {
	lis, _ := net.Listen("tcp", ":50052")
	grpcSrv := grpc.NewServer()
	handler := NewHandler(svc, uc)
	taskpb.RegisterTaskServiceServer(grpcSrv, handler)
	return grpcSrv.Serve(lis)
}
