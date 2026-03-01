package controller

import (
	"context"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/service"

	authpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
)

type OAuthController struct {
	Service *service.OAuthService
}

func NewOAuthController(service *service.OAuthService) *OAuthController {
	return &OAuthController{Service: service}
}

func (c *OAuthController) OAuthRegister(ctx context.Context, req *adminpb.OAuthRegisterRequest) (*adminpb.OAuthRegisterResponse, error) {
	authReq := &authpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
		Birthday: req.Birthday,
	}
	resp, err := c.Service.Register(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return &adminpb.OAuthRegisterResponse{
		Ok:      resp.Ok,
		Message: resp.Message,
	}, nil
}

func (c *OAuthController) OAuthToken(ctx context.Context, req *adminpb.OAuthTokenRequest) (*adminpb.OAuthTokenResponse, error) {
	authReq := &authpb.TokenRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := c.Service.Token(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return &adminpb.OAuthTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (c *OAuthController) OAuthVerify(ctx context.Context, req *adminpb.OAuthVerifyRequest) (*adminpb.OAuthVerifyResponse, error) {
	authReq := &authpb.VerifyRequest{
		AccessToken: req.AccessToken,
	}
	resp, err := c.Service.Verify(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return &adminpb.OAuthVerifyResponse{
		Valid: resp.Valid,
		Email: resp.Email,
		Name:  resp.Name,
		Exp:   resp.Exp,
	}, nil
}

func (c *OAuthController) OAuthRefresh(ctx context.Context, req *adminpb.OAuthRefreshRequest) (*adminpb.OAuthRefreshResponse, error) {
	authReq := &authpb.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
	resp, err := c.Service.Refresh(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return &adminpb.OAuthRefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}
