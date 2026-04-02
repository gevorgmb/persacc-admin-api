package server

import (
	"context"

	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/controller"
	"persacc/internal/service"

	authpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
)

type AdminServer struct {
	adminpb.UnimplementedAdminServiceServer
	DB               *gorm.DB
	AuthClient       authpb.OAuthClient
	UserCtrl         *controller.UserController
	RoleCtrl         *controller.RoleController
	CustomerCtrl     *controller.CustomerController
	PermissionCtrl   *controller.PermissionController
	OAuthCtrl        *controller.OAuthController
	OrganizationCtrl *controller.OrganizationController
	ProductCtrl      *controller.ProductController
}

func NewAdminServer(db *gorm.DB, authClient authpb.OAuthClient) *AdminServer {
	userService := service.NewUserService(db)
	roleService := service.NewRoleService(db)
	customerService := service.NewCustomerService(db)
	permissionService := service.NewPermissionService(db)
	oauthService := service.NewOAuthService(authClient)

	return &AdminServer{
		DB:               db,
		AuthClient:       authClient,
		UserCtrl:         controller.NewUserController(userService),
		RoleCtrl:         controller.NewRoleController(roleService),
		CustomerCtrl:     controller.NewCustomerController(customerService),
		PermissionCtrl:   controller.NewPermissionController(permissionService),
		OAuthCtrl:        controller.NewOAuthController(oauthService),
		OrganizationCtrl: controller.NewOrganizationController(service.NewOrganizationService(db)),
		ProductCtrl:      controller.NewProductController(service.NewProductService(db)),
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

// --- Organization CRUD ---

func (s *AdminServer) CreateOrganization(ctx context.Context, req *adminpb.CreateOrganizationRequest) (*adminpb.CreateOrganizationResponse, error) {
	return s.OrganizationCtrl.Create(ctx, req)
}

func (s *AdminServer) GetOrganization(ctx context.Context, req *adminpb.GetOrganizationRequest) (*adminpb.GetOrganizationResponse, error) {
	return s.OrganizationCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdateOrganization(ctx context.Context, req *adminpb.UpdateOrganizationRequest) (*adminpb.UpdateOrganizationResponse, error) {
	return s.OrganizationCtrl.Update(ctx, req)
}

func (s *AdminServer) DeleteOrganization(ctx context.Context, req *adminpb.DeleteOrganizationRequest) (*adminpb.DeleteOrganizationResponse, error) {
	return s.OrganizationCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListOrganizations(ctx context.Context, req *adminpb.ListOrganizationsRequest) (*adminpb.ListOrganizationsResponse, error) {
	return s.OrganizationCtrl.List(ctx, req)
}

// --- Product CRUD ---

func (s *AdminServer) CreateProduct(ctx context.Context, req *adminpb.CreateProductRequest) (*adminpb.CreateProductResponse, error) {
	return s.ProductCtrl.Create(ctx, req)
}

func (s *AdminServer) GetProduct(ctx context.Context, req *adminpb.GetProductRequest) (*adminpb.GetProductResponse, error) {
	return s.ProductCtrl.Get(ctx, req)
}

func (s *AdminServer) UpdateProduct(ctx context.Context, req *adminpb.UpdateProductRequest) (*adminpb.UpdateProductResponse, error) {
	return s.ProductCtrl.Update(ctx, req)
}

func (s *AdminServer) DeleteProduct(ctx context.Context, req *adminpb.DeleteProductRequest) (*adminpb.DeleteProductResponse, error) {
	return s.ProductCtrl.Delete(ctx, req)
}

func (s *AdminServer) ListProducts(ctx context.Context, req *adminpb.ListProductsRequest) (*adminpb.ListProductsResponse, error) {
	return s.ProductCtrl.List(ctx, req)
}

// --- OAuth Proxy ---

func (s *AdminServer) OAuthRegister(ctx context.Context, req *adminpb.OAuthRegisterRequest) (*adminpb.OAuthRegisterResponse, error) {
	return s.OAuthCtrl.OAuthRegister(ctx, req)
}

func (s *AdminServer) OAuthToken(ctx context.Context, req *adminpb.OAuthTokenRequest) (*adminpb.OAuthTokenResponse, error) {
	return s.OAuthCtrl.OAuthToken(ctx, req)
}

func (s *AdminServer) OAuthVerify(ctx context.Context, req *adminpb.OAuthVerifyRequest) (*adminpb.OAuthVerifyResponse, error) {
	return s.OAuthCtrl.OAuthVerify(ctx, req)
}

func (s *AdminServer) OAuthRefresh(ctx context.Context, req *adminpb.OAuthRefreshRequest) (*adminpb.OAuthRefreshResponse, error) {
	return s.OAuthCtrl.OAuthRefresh(ctx, req)
}
