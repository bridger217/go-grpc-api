package v1

import (
	"context"
	"database/sql"

	v1 "github.com/bridger217/go-grpc-api/pkg/api/v1"
	"github.com/bridger217/go-grpc-api/pkg/db"
	"github.com/bridger217/go-grpc-api/pkg/middleware"
)

type userServiceServer struct {
	userTable *db.UsersTable
	v1.UnimplementedUserServiceServer
}

func NewUserServiceServer(database *sql.DB) *userServiceServer {
	t := db.NewUsersTable(database)
	return &userServiceServer{userTable: t}
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	id, err := middleware.UserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.userTable.InsertUser(id, req.User.Username, req.User.FirstName, req.User.LastName)
	u := req.User
	u.Id = id
	return u, err
}

func (s *userServiceServer) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	u, err := s.userTable.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return u.ToExternal(), nil
}
