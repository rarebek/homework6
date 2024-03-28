package mocktest

import (
	pbu "EXAM3/api-gateway/genproto/user_service"
	"context"

	"google.golang.org/grpc"
)

type UserServiceClientI interface {
	CreateUser(ctx context.Context, in *pbu.User, opts ...grpc.CallOption) (*pbu.User, error)
	GetUserByUsername(ctx context.Context, in *pbu.Username, opts ...grpc.CallOption) (*pbu.User, error)
	GetUserByEmail(ctx context.Context, in *pbu.Email, opts ...grpc.CallOption) (*pbu.User, error)
	UpdateUserById(ctx context.Context, in *pbu.User, opts ...grpc.CallOption) (*pbu.User, error)
	GetUserById(ctx context.Context, in *pbu.UserId, opts ...grpc.CallOption) (*pbu.User, error)
	ListUser(ctx context.Context, in *pbu.GetAllUserRequest, opts ...grpc.CallOption) (*pbu.GetAllUserResponse, error)
	DeleteUser(ctx context.Context, in *pbu.UserId, opts ...grpc.CallOption) (*pbu.Empty, error)
	CheckField(ctx context.Context, in *pbu.CheckFieldRequest, opts ...grpc.CallOption) (*pbu.CheckFieldResponse, error)
}

type UserServiceClient struct {
}

func NewUserServiceClient() UserServiceClientI {
	return &UserServiceClient{}
}

func (c *UserServiceClient) CreateUser(ctx context.Context, in *pbu.User, opts ...grpc.CallOption) (*pbu.User, error) {
	user := pbu.User{
		Id:       "1",
		Name:     "Nodirbek",
		Age:      17,
		Username: "rarebek",
		Email:    "nomonovn2@gmail.com",
	}
	return &user, nil
}

func (c *UserServiceClient) GetUserByUsername(ctx context.Context, in *pbu.Username, opts ...grpc.CallOption) (*pbu.User, error) {
	user := pbu.User{
		Id:       "1",
		Name:     "Nodirbek",
		Age:      17,
		Username: "rarebek",
		Email:    "nomonovn2@gmail.com",
	}
	return &user, nil
}

func (c *UserServiceClient) GetUserByEmail(ctx context.Context, in *pbu.Email, opts ...grpc.CallOption) (*pbu.User, error) {
	user := pbu.User{
		Id:       "dsfhdsjfhl",
		Name:     "Nodirbek",
		Age:      17,
		Username: "rarebek",
		Email:    "nomonovn2@gmail.com",
	}
	return &user, nil
}

func (c *UserServiceClient) UpdateUserById(ctx context.Context, in *pbu.User, opts ...grpc.CallOption) (*pbu.User, error) {
	return in, nil
}

func (c *UserServiceClient) GetUserById(ctx context.Context, in *pbu.UserId, opts ...grpc.CallOption) (*pbu.User, error) {
	user := pbu.User{
		Id:       "1",
		Name:     "Nodirbek",
		Age:      17,
		Username: "rarebek",
		Email:    "nomonovn2@gmail.com",
	}
	return &user, nil
}

func (c *UserServiceClient) ListUser(ctx context.Context, in *pbu.GetAllUserRequest, opts ...grpc.CallOption) (*pbu.GetAllUserResponse, error) {
	user := pbu.User{
		Id:       "dsfhdsjfhl",
		Name:     "Nodirbek",
		Age:      17,
		Username: "rarebek",
		Email:    "nomonovn2@gmail.com",
	}
	resp := pbu.GetAllUserResponse{
		Count: 1,
		Users: []*pbu.User{
			&user,
			&user,
			&user,
		},
	}
	return &resp, nil
}

func (c *UserServiceClient) DeleteUser(ctx context.Context, in *pbu.UserId, opts ...grpc.CallOption) (*pbu.Empty, error) {
	return &pbu.Empty{}, nil
}

func (c *UserServiceClient) CheckField(ctx context.Context, in *pbu.CheckFieldRequest, opts ...grpc.CallOption) (*pbu.CheckFieldResponse, error) {
	return &pbu.CheckFieldResponse{
		Status: false,
	}, nil
}
