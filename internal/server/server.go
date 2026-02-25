package server

import (
	"context"

	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/controller"

	authpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
)

type AdminServer struct {
	adminpb.UnimplementedAdminServiceServer
	DB             *gorm.DB
	AuthClient     authpb.OAuthClient
	UserCtrl       *controller.UserController
	RoleCtrl       *controller.RoleController
	CustomerCtrl   *controller.CustomerController
	PermissionCtrl *controller.PermissionController
}

func NewAdminServer(db *gorm.DB, authClient authpb.OAuthClient) *AdminServer {
	return &AdminServer{
		DB:             db,
		AuthClient:     authClient,
		UserCtrl:       controller.NewUserController(db),
		RoleCtrl:       controller.NewRoleController(db),
		CustomerCtrl:   controller.NewCustomerController(db),
		PermissionCtrl: controller.NewPermissionController(db),
	}
}

// --- User CRUD ---

func (s *AdminServer) Register(ctx context.Context, req *adminpb.RegisterRequest) (*adminpb.RegisterResponse, error) {
	return s.UserCtrl.Register(ctx, req)
}

func (s *AdminServer) CreateUser(ctx context.Context, req *adminpb.CreateUserRequest) (*adminpb.CreateUserResponse, error) {
	return s.UserCtrl.Create(ctx, req)
}

func (s *AdminServer) GetUser(ctx context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	return s.UserCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdateUser(ctx context.Context, req *adminpb.UpdateUserRequest) (*adminpb.UpdateUserResponse, error) {
	return s.UserCtrl.Update(ctx, req)
}

func (s *AdminServer) DeleteUser(ctx context.Context, req *adminpb.DeleteUserRequest) (*adminpb.DeleteUserResponse, error) {
	return s.UserCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListUsers(ctx context.Context, req *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	return s.UserCtrl.List(ctx, req)
}

// --- Role CRUD ---

func (s *AdminServer) CreateRole(ctx context.Context, req *adminpb.CreateRoleRequest) (*adminpb.CreateRoleResponse, error) {
	return s.RoleCtrl.Create(ctx, req)
}

func (s *AdminServer) GetRole(ctx context.Context, req *adminpb.GetRoleRequest) (*adminpb.GetRoleResponse, error) {
	return s.RoleCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdateRole(ctx context.Context, req *adminpb.UpdateRoleRequest) (*adminpb.UpdateRoleResponse, error) {
	return s.RoleCtrl.Update(ctx, req)
}

func (s *AdminServer) DeleteRole(ctx context.Context, req *adminpb.DeleteRoleRequest) (*adminpb.DeleteRoleResponse, error) {
	return s.RoleCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListRoles(ctx context.Context, req *adminpb.ListRolesRequest) (*adminpb.ListRolesResponse, error) {
	return s.RoleCtrl.List(ctx, req)
}

// --- Customer CRUD ---

func (s *AdminServer) CreateCustomer(ctx context.Context, req *adminpb.CreateCustomerRequest) (*adminpb.CreateCustomerResponse, error) {
	return s.CustomerCtrl.Create(ctx, req)
}

func (s *AdminServer) GetCustomer(ctx context.Context, req *adminpb.GetCustomerRequest) (*adminpb.GetCustomerResponse, error) {
	return s.CustomerCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdateCustomer(ctx context.Context, req *adminpb.UpdateCustomerRequest) (*adminpb.UpdateCustomerResponse, error) {
	return s.CustomerCtrl.Update(ctx, req)
}

func (s *AdminServer) DeleteCustomer(ctx context.Context, req *adminpb.DeleteCustomerRequest) (*adminpb.DeleteCustomerResponse, error) {
	return s.CustomerCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListCustomers(ctx context.Context, req *adminpb.ListCustomersRequest) (*adminpb.ListCustomersResponse, error) {
	return s.CustomerCtrl.List(ctx, req)
}

// --- Permission CRUD ---

func (s *AdminServer) CreatePermission(ctx context.Context, req *adminpb.CreatePermissionRequest) (*adminpb.CreatePermissionResponse, error) {
	return s.PermissionCtrl.Create(ctx, req)
}

func (s *AdminServer) GetPermission(ctx context.Context, req *adminpb.GetPermissionRequest) (*adminpb.GetPermissionResponse, error) {
	return s.PermissionCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdatePermission(ctx context.Context, req *adminpb.UpdatePermissionRequest) (*adminpb.UpdatePermissionResponse, error) {
	return s.PermissionCtrl.Update(ctx, req)
}

func (s *AdminServer) DeletePermission(ctx context.Context, req *adminpb.DeletePermissionRequest) (*adminpb.DeletePermissionResponse, error) {
	return s.PermissionCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListPermissions(ctx context.Context, req *adminpb.ListPermissionsRequest) (*adminpb.ListPermissionsResponse, error) {
	return s.PermissionCtrl.List(ctx, req)
}
