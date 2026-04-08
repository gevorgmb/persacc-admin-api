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

type ProductCategoryController struct {
	Service *service.ProductCategoryService
}

func NewProductCategoryController(service *service.ProductCategoryService) *ProductCategoryController {
	return &ProductCategoryController{Service: service}
}

func (c *ProductCategoryController) Create(ctx context.Context, req *adminpb.CreateProductCategoryRequest) (*adminpb.CreateProductCategoryResponse, error) {
	orgId := ctx.Value("organization_id").(int64)

	category := entity.ProductCategory{
		OrganizationID: orgId,
		Name:           req.Name,
	}

	if req.Description != "" {
		desc := req.Description
		category.Description = &desc
	}

	if err := c.Service.Create(ctx, &category); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product category: %v", err)
	}

	return &adminpb.CreateProductCategoryResponse{
		Category: ConvertProductCategoryToProto(category),
	}, nil
}

func (c *ProductCategoryController) Get(ctx context.Context, req *adminpb.GetProductCategoryRequest) (*adminpb.GetProductCategoryResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	category, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product category: %v", err)
	}

	return &adminpb.GetProductCategoryResponse{
		Category: ConvertProductCategoryToProto(*category),
	}, nil
}

func (c *ProductCategoryController) Update(ctx context.Context, req *adminpb.UpdateProductCategoryRequest) (*adminpb.UpdateProductCategoryResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	category, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find product category: %v", err)
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		desc := req.Description
		category.Description = &desc
	}

	if err := c.Service.Update(ctx, category, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product category: %v", err)
	}

	return &adminpb.UpdateProductCategoryResponse{
		Category: ConvertProductCategoryToProto(*category),
	}, nil
}

func (c *ProductCategoryController) Delete(ctx context.Context, req *adminpb.DeleteProductCategoryRequest) (*adminpb.DeleteProductCategoryResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	if err := c.Service.Delete(ctx, req.Id, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product category: %v", err)
	}
	return &adminpb.DeleteProductCategoryResponse{Success: true}, nil
}

func (c *ProductCategoryController) List(ctx context.Context, req *adminpb.ListProductCategoriesRequest) (*adminpb.ListProductCategoriesResponse, error) {
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

	categories, total, err := c.Service.List(ctx, limit, offset, orgId, filters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list product categories: %v", err)
	}

	var protoCategories []*adminpb.ProductCategory
	for _, cat := range categories {
		protoCategories = append(protoCategories, ConvertProductCategoryToProto(cat))
	}

	return &adminpb.ListProductCategoriesResponse{
		Categories: protoCategories,
		Total:      int32(total),
		Page:       int32(page),
		Limit:      int32(limit),
	}, nil
}

func ConvertProductCategoryToProto(cat entity.ProductCategory) *adminpb.ProductCategory {
	var description string
	if cat.Description != nil {
		description = *cat.Description
	}

	return &adminpb.ProductCategory{
		Id:             cat.ID,
		OrganizationId: cat.OrganizationID,
		Name:           cat.Name,
		Description:    description,
		CreatedAt:      cat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      cat.UpdatedAt.Format(time.RFC3339),
	}
}
