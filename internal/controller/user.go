package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
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

	// Check if user already exists
	var existingUser entity.User
	if err := c.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, "failed to check existing user: %v", err)
	}

	// Find the 'user' role
	var role entity.Role
	if err := c.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Internal, "'user' role not found in the system")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve 'user' role: %v", err)
	}

	user := entity.User{
		Name:   name,
		Email:  email,
		RoleID: role.ID,
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &adminpb.RegisterResponse{
		User: ConvertUserToProto(user),
	}, nil
}

func (c *UserController) Create(ctx context.Context, req *adminpb.CreateUserRequest) (*adminpb.CreateUserResponse, error) {
	user := entity.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleId,
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &adminpb.CreateUserResponse{
		User: ConvertUserToProto(user),
	}, nil
}

func (c *UserController) Get(ctx context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	var user entity.User
	if err := c.DB.First(&user, "id = ?", req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &adminpb.GetUserResponse{
		User: ConvertUserToProto(user),
	}, nil
}

func (c *UserController) Update(ctx context.Context, req *adminpb.UpdateUserRequest) (*adminpb.UpdateUserResponse, error) {
	var user entity.User
	if err := c.DB.First(&user, "id = ?", req.Id).Error; err != nil {
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

	if err := c.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &adminpb.UpdateUserResponse{
		User: ConvertUserToProto(user),
	}, nil
}

func (c *UserController) Delete(ctx context.Context, req *adminpb.DeleteUserRequest) (*adminpb.DeleteUserResponse, error) {
	if err := c.DB.Delete(&entity.User{}, "id = ?", req.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &adminpb.DeleteUserResponse{Success: true}, nil
}

func (c *UserController) List(ctx context.Context, req *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	var users []entity.User
	var total int64

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	c.DB.Model(&entity.User{}).Count(&total)
	if err := c.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
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
