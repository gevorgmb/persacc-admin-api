package service

import (
	"context"

	authpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
)

type OAuthService struct {
	Client authpb.OAuthClient
}

func NewOAuthService(client authpb.OAuthClient) *OAuthService {
	return &OAuthService{Client: client}
}

func (s *OAuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return s.Client.Register(ctx, req)
}

func (s *OAuthService) Token(ctx context.Context, req *authpb.TokenRequest) (*authpb.TokenResponse, error) {
	return s.Client.Token(ctx, req)
}

func (s *OAuthService) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	return s.Client.Verify(ctx, req)
}

func (s *OAuthService) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.TokenResponse, error) {
	return s.Client.Refresh(ctx, req)
}
