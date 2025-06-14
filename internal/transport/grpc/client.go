package grpc

import (
	userpb "github.com/4nar1k/project-protos/proto/user"
	"google.golang.org/grpc"
)

func NewUserClient(addr string) (userpb.UserServiceClient, *grpc.ClientConn, error) {
	// Устанавливаем соединение с Users-сервисом
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	// Создаем клиента для Users-сервиса
	client := userpb.NewUserServiceClient(conn)

	return client, conn, nil
}
