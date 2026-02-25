package server

import (
	"context"
	"errors"
	"log"
	"strings"

	"persacc/internal/entity"

	oauthpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
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

		// Skip auth for reflection
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
			return handler(ctx, req)
		}
		log.Println("Not reflection")

		// 1. Extract Authorization from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}
		log.Println("Metadata is provided")

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

		// 4. User is authenticated, now check authorization (Role)
		userEmail := verifyResp.Email
		var user entity.User
		// Preload Role to check its name
		if err := i.DB.Preload("Role").First(&user, "email = ?", userEmail).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Errorf(codes.PermissionDenied, "user not found in admin system")
			}
			return nil, status.Errorf(codes.Internal, "failed to check user role: %v", err)
		}

		// 5. Check if role is "admin"
		if user.Role.Name != "admin" {
			return nil, status.Errorf(codes.PermissionDenied, "access denied: requires admin role")
		}

		// 5. Proceed
		return handler(ctx, req)
	}
}
