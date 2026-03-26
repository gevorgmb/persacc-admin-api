package server

import (
	"context"
	"errors"
	"log"
	"strings"

	"persacc/internal/entity"

	oauthpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthInterceptor struct {
	DB         *gorm.DB
	AuthClient oauthpb.OAuthClient
}

func NewAuthInterceptor(db *gorm.DB, authClient oauthpb.OAuthClient) *AuthInterceptor {
	return &AuthInterceptor{
		DB:         db,
		AuthClient: authClient,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Log the method being called for debugging
		log.Printf("Checking access for method: %s", info.FullMethod)

		// Skip auth for reflection and OAuth proxy methods
		publicMethods := map[string]bool{
			"/admin.AdminService/OAuthRegister": true,
			"/admin.AdminService/OAuthToken":    true,
			"/admin.AdminService/OAuthVerify":   true,
			"/admin.AdminService/OAuthRefresh":  true,
		}

		if strings.HasPrefix(info.FullMethod, "/grpc.reflection") || publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 1. Extract Authorization from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}
		log.Println("Metadata is provided:")
		for k, v := range md {
			log.Printf("  Metadata: %q: %v", k, v)
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
		log.Println("Authorization token is provided")

		accessToken := values[0]
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		// 2. Call OAuth Verify
		verifyResp, err := i.AuthClient.Verify(ctx, &oauthpb.VerifyRequest{
			AccessToken: accessToken,
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "failed to verify token: %v", err)
		}
		log.Println("Token is valid")

		if !verifyResp.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "token is invalid")
		}
		log.Println("Token is valid")

		// 3. Skip DB check if this is the Register method
		if info.FullMethod == "/admin.AdminService/Register" {
			// Pass email and name in context to the Register handler
			ctx = context.WithValue(ctx, "email", verifyResp.Email)
			ctx = context.WithValue(ctx, "name", verifyResp.Name)
			return handler(ctx, req)
		}

		// 4. User is authenticated, now sync with local DB
		userEmail := verifyResp.Email
		userUuid := verifyResp.GetUuid()
		var user entity.User

		foundByUuid := false
		if userUuid != "" {
			if err := i.DB.Preload("Role").First(&user, "uuid = ?", userUuid).Error; err == nil {
				foundByUuid = true
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Errorf(codes.Internal, "failed to check user by uuid: %v", err)
			}
		}

		if !foundByUuid {
			if err := i.DB.Preload("Role").First(&user, "email = ?", userEmail).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// User not found, create a new one
					user = entity.User{
						Name:   verifyResp.Name,
						Email:  userEmail,
						Uuid:   userUuid,
						RoleID: 1, // Default role
					}
					if err := i.DB.Create(&user).Error; err != nil {
						return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
					}
					// Reload to get the Role for the role checks below
					if err := i.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
						return nil, status.Errorf(codes.Internal, "failed to load created user: %v", err)
					}
				} else {
					return nil, status.Errorf(codes.Internal, "failed to check user by email: %v", err)
				}
			} else if user.Uuid == "" && userUuid != "" {
				// User found by email but missing UUID, update it
				user.Uuid = userUuid
				if err := i.DB.Save(&user).Error; err != nil {
					return nil, status.Errorf(codes.Internal, "failed to update user uuid: %v", err)
				}
			}
		}

		// 5. Define paths that explicitly require admin access
		requireAdminPaths := map[string]bool{
			"/admin.AdminService/CreateUser":       true,
			"/admin.AdminService/UpdateUser":       true,
			"/admin.AdminService/DeleteUser":       true,
			"/admin.AdminService/CreateRole":       true,
			"/admin.AdminService/UpdateRole":       true,
			"/admin.AdminService/DeleteRole":       true,
			"/admin.AdminService/CreatePermission": true,
			"/admin.AdminService/UpdatePermission": true,
			"/admin.AdminService/DeletePermission": true,
		}

		// 6. Check if this method requires admin and if the user is an admin
		if requireAdminPaths[info.FullMethod] {
			if user.Role.Name != "admin" {
				return nil, status.Errorf(codes.PermissionDenied, "access denied: method requires admin role")
			}
		}

		// 8. Enforce Organization ID for all requests except exempted ones
		isExempted := strings.Contains(info.FullMethod, "OAuth") ||
			strings.Contains(info.FullMethod, "Organization") ||
			info.FullMethod == "/admin.AdminService/Register" ||
			strings.HasPrefix(info.FullMethod, "/grpc.reflection")

		if !isExempted {
			orgId := ""
			if vals := md.Get("organization_id"); len(vals) > 0 {
				orgId = vals[0]
			}
			if orgId == "" {
				if vals := md.Get("organization-id"); len(vals) > 0 {
					orgId = vals[0]
				}
			}

			if strings.TrimSpace(orgId) == "" {
				st := status.New(codes.InvalidArgument, "missing organization_id header")
				v := &errdetails.BadRequest_FieldViolation{
					Field:       "organization_id",
					Description: "The organization_id header is required for this request",
				}
				br := &errdetails.BadRequest{}
				br.FieldViolations = append(br.FieldViolations, v)
				st, err := st.WithDetails(br)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "failed to attach error details: %v", err)
				}
				return nil, st.Err()
			}
		}

		// 9. Proceed
		ctx = context.WithValue(ctx, "user_id", user.ID)
		return handler(ctx, req)
	}
}
