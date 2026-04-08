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

type ProductController struct {
	Service *service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{Service: service}
}

func (c *ProductController) Create(ctx context.Context, req *adminpb.CreateProductRequest) (*adminpb.CreateProductResponse, error) {
	orgId := ctx.Value("organization_id").(int64)

	product := entity.Product{
		OrganizationID: orgId,
		SKU:            req.Sku,
		Name:           req.Name,
	}

	if req.Description != "" {
		desc := req.Description
		product.Description = &desc
	}

	if len(req.AdditionalDetails) > 0 {
		product.ProductDetails = &entity.ProductDetail{
			AdditionalDetails: make(map[string]interface{}),
		}
		for k, v := range req.AdditionalDetails {
			product.ProductDetails.AdditionalDetails[k] = v
		}
	}

	if err := c.Service.Create(ctx, &product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &adminpb.CreateProductResponse{
		Product: ConvertProductToProto(product),
	}, nil
}

func (c *ProductController) Get(ctx context.Context, req *adminpb.GetProductRequest) (*adminpb.GetProductResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	product, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	return &adminpb.GetProductResponse{
		Product: ConvertProductToProto(*product),
	}, nil
}

func (c *ProductController) Update(ctx context.Context, req *adminpb.UpdateProductRequest) (*adminpb.UpdateProductResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	product, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find product: %v", err)
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Sku != "" {
		product.SKU = req.Sku
	}
	if req.Description != "" {
		desc := req.Description
		product.Description = &desc
	}

	if len(req.AdditionalDetails) > 0 {
		if product.ProductDetails == nil {
			product.ProductDetails = &entity.ProductDetail{
				AdditionalDetails: make(map[string]interface{}),
			}
		} else if product.ProductDetails.AdditionalDetails == nil {
			product.ProductDetails.AdditionalDetails = make(map[string]interface{})
		}
		for k, v := range req.AdditionalDetails {
			product.ProductDetails.AdditionalDetails[k] = v
		}
	}

	if err := c.Service.Update(ctx, product, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &adminpb.UpdateProductResponse{
		Product: ConvertProductToProto(*product),
	}, nil
}

func (c *ProductController) Delete(ctx context.Context, req *adminpb.DeleteProductRequest) (*adminpb.DeleteProductResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	if err := c.Service.Delete(ctx, req.Id, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &adminpb.DeleteProductResponse{Success: true}, nil
}

func (c *ProductController) List(ctx context.Context, req *adminpb.ListProductsRequest) (*adminpb.ListProductsResponse, error) {
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
	if req.Sku != "" {
		filters["sku"] = req.Sku
	}
	if req.Description != "" {
		filters["description"] = req.Description
	}

	products, total, err := c.Service.List(ctx, limit, offset, orgId, filters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	var protoProducts []*adminpb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, ConvertProductToProto(p))
	}

	return &adminpb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Page:     int32(page),
		Limit:    int32(limit),
	}, nil
}

func ConvertProductToProto(p entity.Product) *adminpb.Product {
	var description string
	if p.Description != nil {
		description = *p.Description
	}

	additionalDetails := make(map[string]string)
	if p.ProductDetails != nil && p.ProductDetails.AdditionalDetails != nil {
		for k, v := range p.ProductDetails.AdditionalDetails {
			if strVal, ok := v.(string); ok {
				additionalDetails[k] = strVal
			}
		}
	}

	return &adminpb.Product{
		Id:                p.ID,
		OrganizationId:    p.OrganizationID,
		Sku:               p.SKU,
		Name:              p.Name,
		Description:       description,
		CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         p.UpdatedAt.Format(time.RFC3339),
		DeletedAt:         p.DeletedAt.Time.Format(time.RFC3339),
		AdditionalDetails: additionalDetails,
	}
}
