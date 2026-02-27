package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
	"persacc/internal/service"
)

type UserController struct {
	Service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{Service: service}
}

func (c *UserController) Register(ctx context.Context, req *adminpb.RegisterRequest) (*adminpb.RegisterResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, status.Errorf(codes.Internal, "email not found in context")
	}

	name, ok := ctx.Value("name").(string)
	if !ok || name == "" {
		return nil, status.Errorf(codes.Internal, "name not found in context")
	}

	user, err := c.Service.Register(ctx, email, name)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return nil, status.Errorf(codes.AlreadyExists, "%v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &adminpb.RegisterResponse{
		User: ConvertUserToProto(*user),
	}, nil
}

func (c *UserController) Create(ctx context.Context, req *adminpb.CreateUserRequest) (*adminpb.CreateUserResponse, error) {
	user := entity.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleId,
	}

	if err := c.Service.Create(ctx, &user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &adminpb.CreateUserResponse{
		User: ConvertUserToProto(user),
	}, nil
}

func (c *UserController) Get(ctx context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	user, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &adminpb.GetUserResponse{
		User: ConvertUserToProto(*user),
	}, nil
}

func (c *UserController) Update(ctx context.Context, req *adminpb.UpdateUserRequest) (*adminpb.UpdateUserResponse, error) {
	user, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RoleId != 0 {
		user.RoleID = req.RoleId
	}

	if err := c.Service.Update(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &adminpb.UpdateUserResponse{
		User: ConvertUserToProto(*user),
	}, nil
}

func (c *UserController) Delete(ctx context.Context, req *adminpb.DeleteUserRequest) (*adminpb.DeleteUserResponse, error) {
	if err := c.Service.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &adminpb.DeleteUserResponse{Success: true}, nil
}

func (c *UserController) List(ctx context.Context, req *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	users, total, err := c.Service.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var protoUsers []*adminpb.User
	for _, u := range users {
		protoUsers = append(protoUsers, ConvertUserToProto(u))
	}

	return &adminpb.ListUsersResponse{
		Users: protoUsers,
		Total: int32(total),
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}

func ConvertUserToProto(u entity.User) *adminpb.User {
	return &adminpb.User{
		Id:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		RoleId: u.RoleID,
	}
}
