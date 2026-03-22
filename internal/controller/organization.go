package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
	"persacc/internal/service"
)

type OrganizationController struct {
	Service *service.OrganizationService
}

func NewOrganizationController(service *service.OrganizationService) *OrganizationController {
	return &OrganizationController{Service: service}
}

func (c *OrganizationController) Create(ctx context.Context, req *adminpb.CreateOrganizationRequest) (*adminpb.CreateOrganizationResponse, error) {
	org := entity.Organization{
		OwnerID:     req.OwnerId,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := c.Service.Create(ctx, &org); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create organization: %v", err)
	}

	return &adminpb.CreateOrganizationResponse{
		Organization: ConvertOrganizationToProto(org),
	}, nil
}

func (c *OrganizationController) Get(ctx context.Context, req *adminpb.GetOrganizationRequest) (*adminpb.GetOrganizationResponse, error) {
	org, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &adminpb.GetOrganizationResponse{
				Message: "organization does not added, add it now",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get organization: %v", err)
	}

	return &adminpb.GetOrganizationResponse{
		Organization: ConvertOrganizationToProto(*org),
	}, nil
}

func (c *OrganizationController) Update(ctx context.Context, req *adminpb.UpdateOrganizationRequest) (*adminpb.UpdateOrganizationResponse, error) {
	org, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find organization: %v", err)
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Description != "" {
		org.Description = req.Description
	}

	if err := c.Service.Update(ctx, org); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update organization: %v", err)
	}

	return &adminpb.UpdateOrganizationResponse{
		Organization: ConvertOrganizationToProto(*org),
	}, nil
}

func (c *OrganizationController) Delete(ctx context.Context, req *adminpb.DeleteOrganizationRequest) (*adminpb.DeleteOrganizationResponse, error) {
	if err := c.Service.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete organization: %v", err)
	}
	return &adminpb.DeleteOrganizationResponse{Success: true}, nil
}

func (c *OrganizationController) List(ctx context.Context, req *adminpb.ListOrganizationsRequest) (*adminpb.ListOrganizationsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	orgs, total, err := c.Service.List(ctx, limit, offset, ctx.Value("user_id").(int64))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list organizations: %v", err)
	}

	var protoOrgs []*adminpb.Organization
	for _, o := range orgs {
		protoOrgs = append(protoOrgs, ConvertOrganizationToProto(o))
	}

	return &adminpb.ListOrganizationsResponse{
		Organizations: protoOrgs,
		Total:         int32(total),
		Page:          int32(page),
		Limit:         int32(limit),
	}, nil
}

func ConvertOrganizationToProto(o entity.Organization) *adminpb.Organization {
	return &adminpb.Organization{
		Id:          o.ID,
		OwnerId:     o.OwnerID,
		Name:        o.Name,
		Description: o.Description,
		CreatedAt:   timestamppb.New(o.CreatedAt),
		UpdatedAt:   timestamppb.New(o.UpdatedAt),
	}
}
