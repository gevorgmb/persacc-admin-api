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

type RoleController struct {
	Service *service.RoleService
}

func NewRoleController(service *service.RoleService) *RoleController {
	return &RoleController{Service: service}
}

func (c *RoleController) Create(ctx context.Context, req *adminpb.CreateRoleRequest) (*adminpb.CreateRoleResponse, error) {
	role := entity.Role{
		Name: req.Name,
	}

	if err := c.Service.Create(ctx, &role, req.PermissionIds); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create role: %v", err)
	}

	return &adminpb.CreateRoleResponse{
		Role: ConvertRoleToProto(role),
	}, nil
}

func (c *RoleController) Get(ctx context.Context, req *adminpb.GetRoleRequest) (*adminpb.GetRoleResponse, error) {
	role, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "role not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get role: %v", err)
	}

	return &adminpb.GetRoleResponse{
		Role: ConvertRoleToProto(*role),
	}, nil
}

func (c *RoleController) Update(ctx context.Context, req *adminpb.UpdateRoleRequest) (*adminpb.UpdateRoleResponse, error) {
	role, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "role not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find role: %v", err)
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	updatePerms := req.PermissionIds != nil
	if err := c.Service.Update(ctx, role, req.PermissionIds, updatePerms); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	return &adminpb.UpdateRoleResponse{
		Role: ConvertRoleToProto(*role),
	}, nil
}

func (c *RoleController) Delete(ctx context.Context, req *adminpb.DeleteRoleRequest) (*adminpb.DeleteRoleResponse, error) {
	if err := c.Service.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete role: %v", err)
	}
	return &adminpb.DeleteRoleResponse{Success: true}, nil
}

func (c *RoleController) List(ctx context.Context, req *adminpb.ListRolesRequest) (*adminpb.ListRolesResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	roles, total, err := c.Service.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}

	var protoRoles []*adminpb.Role
	for _, r := range roles {
		protoRoles = append(protoRoles, ConvertRoleToProto(r))
	}

	return &adminpb.ListRolesResponse{
		Roles: protoRoles,
		Total: int32(total),
		Page:  int32(page),
		Limit: int32(limit),
	}, nil
}

func ConvertRoleToProto(r entity.Role) *adminpb.Role {
	var protoPerms []*adminpb.Permission
	for _, p := range r.Permissions {
		protoPerms = append(protoPerms, ConvertPermissionToProto(p))
	}

	return &adminpb.Role{
		Id:          r.ID,
		Name:        r.Name,
		Permissions: protoPerms,
	}
}
