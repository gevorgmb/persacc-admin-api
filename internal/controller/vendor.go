package controller

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
	"persacc/internal/service"
)

type VendorController struct {
	Service *service.VendorService
}

func NewVendorController(service *service.VendorService) *VendorController {
	return &VendorController{Service: service}
}

func (c *VendorController) Create(ctx context.Context, req *adminpb.CreateVendorRequest) (*adminpb.CreateVendorResponse, error) {
	vendor := entity.Vendor{
		Name: req.Name,
	}

	if req.Domain != "" {
		vendor.Domain = &req.Domain
	}
	if req.Description != "" {
		vendor.Description = &req.Description
	}

	if err := c.Service.Create(ctx, &vendor); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create vendor: %v", err)
	}

	return &adminpb.CreateVendorResponse{
		Vendor: ConvertVendorToProto(vendor),
	}, nil
}

func (c *VendorController) Get(ctx context.Context, req *adminpb.GetVendorRequest) (*adminpb.GetVendorResponse, error) {
	vendor, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "vendor not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get vendor: %v", err)
	}

	return &adminpb.GetVendorResponse{
		Vendor: ConvertVendorToProto(*vendor),
	}, nil
}

func (c *VendorController) Update(ctx context.Context, req *adminpb.UpdateVendorRequest) (*adminpb.UpdateVendorResponse, error) {
	vendor, err := c.Service.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "vendor not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find vendor: %v", err)
	}

	if req.Name != "" {
		vendor.Name = req.Name
	}
	if req.Domain != "" {
		vendor.Domain = &req.Domain
	}
	if req.Description != "" {
		vendor.Description = &req.Description
	}

	if err := c.Service.Update(ctx, vendor); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update vendor: %v", err)
	}

	return &adminpb.UpdateVendorResponse{
		Vendor: ConvertVendorToProto(*vendor),
	}, nil
}

func (c *VendorController) Delete(ctx context.Context, req *adminpb.DeleteVendorRequest) (*adminpb.DeleteVendorResponse, error) {
	if err := c.Service.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete vendor: %v", err)
	}
	return &adminpb.DeleteVendorResponse{Success: true}, nil
}

func (c *VendorController) List(ctx context.Context, req *adminpb.ListVendorsRequest) (*adminpb.ListVendorsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	filters := make(map[string]string)
	if req.Name != "" {
		filters["name"] = req.Name
	}

	vendors, total, err := c.Service.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vendors: %v", err)
	}

	var protoVendors []*adminpb.Vendor
	for _, v := range vendors {
		protoVendors = append(protoVendors, ConvertVendorToProto(v))
	}

	return &adminpb.ListVendorsResponse{
		Vendors: protoVendors,
		Total:   int32(total),
		Page:    int32(page),
		Limit:   int32(limit),
	}, nil
}

func ConvertVendorToProto(v entity.Vendor) *adminpb.Vendor {
	var domain, description string
	if v.Domain != nil {
		domain = *v.Domain
	}
	if v.Description != nil {
		description = *v.Description
	}

	return &adminpb.Vendor{
		Id:          v.ID,
		Name:        v.Name,
		Domain:      domain,
		Description: description,
		CreatedAt:   v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   v.UpdatedAt.Format(time.RFC3339),
	}
}
