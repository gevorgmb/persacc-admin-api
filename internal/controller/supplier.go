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

type SupplierController struct {
	Service *service.SupplierService
}

func NewSupplierController(service *service.SupplierService) *SupplierController {
	return &SupplierController{Service: service}
}

func (c *SupplierController) Create(ctx context.Context, req *adminpb.CreateSupplierRequest) (*adminpb.CreateSupplierResponse, error) {
	orgId := ctx.Value("organization_id").(int64)

	supplier := entity.Supplier{
		Name:           req.Name,
		OrganizationID: &orgId,
	}

	if req.Domain != "" {
		supplier.Domain = &req.Domain
	}
	if req.Phone != "" {
		supplier.Phone = &req.Phone
	}
	if req.Description != "" {
		supplier.Description = &req.Description
	}

	if err := c.Service.Create(ctx, &supplier); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create supplier: %v", err)
	}

	return &adminpb.CreateSupplierResponse{
		Supplier: ConvertSupplierToProto(supplier),
	}, nil
}

func (c *SupplierController) Get(ctx context.Context, req *adminpb.GetSupplierRequest) (*adminpb.GetSupplierResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	supplier, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "supplier not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get supplier: %v", err)
	}

	return &adminpb.GetSupplierResponse{
		Supplier: ConvertSupplierToProto(*supplier),
	}, nil
}

func (c *SupplierController) Update(ctx context.Context, req *adminpb.UpdateSupplierRequest) (*adminpb.UpdateSupplierResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	supplier, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "supplier not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find supplier: %v", err)
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.Domain != "" {
		supplier.Domain = &req.Domain
	}
	if req.Phone != "" {
		supplier.Phone = &req.Phone
	}
	if req.Description != "" {
		supplier.Description = &req.Description
	}

	if err := c.Service.Update(ctx, supplier, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update supplier: %v", err)
	}

	return &adminpb.UpdateSupplierResponse{
		Supplier: ConvertSupplierToProto(*supplier),
	}, nil
}

func (c *SupplierController) Delete(ctx context.Context, req *adminpb.DeleteSupplierRequest) (*adminpb.DeleteSupplierResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	if err := c.Service.Delete(ctx, req.Id, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete supplier: %v", err)
	}
	return &adminpb.DeleteSupplierResponse{Success: true}, nil
}

func (c *SupplierController) List(ctx context.Context, req *adminpb.ListSuppliersRequest) (*adminpb.ListSuppliersResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	orgId := ctx.Value("organization_id").(int64)

	filters := make(map[string]string)
	if req.Name != "" {
		filters["name"] = req.Name
	}

	suppliers, total, err := c.Service.List(ctx, limit, offset, orgId, filters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list suppliers: %v", err)
	}

	var protoSuppliers []*adminpb.Supplier
	for _, s := range suppliers {
		protoSuppliers = append(protoSuppliers, ConvertSupplierToProto(s))
	}

	return &adminpb.ListSuppliersResponse{
		Suppliers: protoSuppliers,
		Total:     int32(total),
		Page:      int32(page),
		Limit:     int32(limit),
	}, nil
}

func ConvertSupplierToProto(s entity.Supplier) *adminpb.Supplier {
	var domain, phone, description string
	if s.Domain != nil {
		domain = *s.Domain
	}
	if s.Phone != nil {
		phone = *s.Phone
	}
	if s.Description != nil {
		description = *s.Description
	}

	return &adminpb.Supplier{
		Id:          s.ID,
		Name:        s.Name,
		Domain:      domain,
		Phone:       phone,
		Description: description,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
}
