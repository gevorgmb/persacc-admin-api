package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/entity"
	"persacc/internal/service"
)

type CustomerController struct {
	Service *service.CustomerService
}

func NewCustomerController(service *service.CustomerService) *CustomerController {
	return &CustomerController{Service: service}
}

func (c *CustomerController) Create(ctx context.Context, req *adminpb.CreateCustomerRequest) (*adminpb.CreateCustomerResponse, error) {
	customer := entity.Customer{
		Name:       req.Name,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Prefix:     req.Prefix,
		MiddleName: req.MiddleName,
		Suffix:     req.Suffix,
		Phone:      req.Phone,
		Email:      req.Email,
	}

	if req.Birthday != "" {
		if t, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			customer.Birthday = &t
		}
	}

	if len(req.AdditionalInfo) > 0 {
		info := make(map[string]interface{}, len(req.AdditionalInfo))
		for k, v := range req.AdditionalInfo {
			info[k] = v
		}
		customer.AdditionalInfo = info
	}

	if req.UserId != 0 {
		uid := req.UserId
		customer.UserID = &uid
	}

	orgId := ctx.Value("organization_id").(int64)

	if err := c.Service.Create(ctx, &customer, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create customer: %v", err)
	}

	return &adminpb.CreateCustomerResponse{
		Customer: ConvertCustomerToProto(customer),
	}, nil
}

func (c *CustomerController) Get(ctx context.Context, req *adminpb.GetCustomerRequest) (*adminpb.GetCustomerResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	customer, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get customer: %v", err)
	}

	return &adminpb.GetCustomerResponse{
		Customer: ConvertCustomerToProto(*customer),
	}, nil
}

func (c *CustomerController) Update(ctx context.Context, req *adminpb.UpdateCustomerRequest) (*adminpb.UpdateCustomerResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	customer, err := c.Service.Get(ctx, req.Id, orgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find customer: %v", err)
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.FirstName != "" {
		customer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		customer.LastName = req.LastName
	}
	if req.Prefix != "" {
		customer.Prefix = req.Prefix
	}
	if req.MiddleName != "" {
		customer.MiddleName = req.MiddleName
	}
	if req.Suffix != "" {
		customer.Suffix = req.Suffix
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Birthday != "" {
		if t, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			customer.Birthday = &t
		}
	}
	if len(req.AdditionalInfo) > 0 {
		info := make(map[string]interface{}, len(req.AdditionalInfo))
		for k, v := range req.AdditionalInfo {
			info[k] = v
		}
		customer.AdditionalInfo = info
	}

	if err := c.Service.Update(ctx, customer, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update customer: %v", err)
	}

	return &adminpb.UpdateCustomerResponse{
		Customer: ConvertCustomerToProto(*customer),
	}, nil
}

func (c *CustomerController) Delete(ctx context.Context, req *adminpb.DeleteCustomerRequest) (*adminpb.DeleteCustomerResponse, error) {
	orgId := ctx.Value("organization_id").(int64)
	if err := c.Service.Delete(ctx, req.Id, orgId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete customer: %v", err)
	}
	return &adminpb.DeleteCustomerResponse{Success: true}, nil
}

func (c *CustomerController) List(ctx context.Context, req *adminpb.ListCustomersRequest) (*adminpb.ListCustomersResponse, error) {
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
	if req.Email != "" {
		filters["email"] = req.Email
	}
	if req.Phone != "" {
		filters["phone"] = req.Phone
	}
	if req.AdditionalInfo != "" {
		filters["additional_info"] = req.AdditionalInfo
	}

	customers, total, err := c.Service.List(ctx, limit, offset, orgId, filters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list customers: %v", err)
	}

	var protoCustomers []*adminpb.Customer
	for _, cu := range customers {
		protoCustomers = append(protoCustomers, ConvertCustomerToProto(cu))
	}

	return &adminpb.ListCustomersResponse{
		Customers: protoCustomers,
		Total:     int32(total),
		Page:      int32(page),
		Limit:     int32(limit),
	}, nil
}

func ConvertCustomerToProto(c entity.Customer) *adminpb.Customer {
	var birthday string
	if c.Birthday != nil {
		birthday = c.Birthday.Format("2006-01-02")
	}

	additionalInfo := make(map[string]string, len(c.AdditionalInfo))
	for k, v := range c.AdditionalInfo {
		additionalInfo[k] = fmt.Sprintf("%v", v)
	}

	var userId int64
	if c.UserID != nil {
		userId = *c.UserID
	}

	return &adminpb.Customer{
		Id:             c.ID,
		Name:           c.Name,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		Prefix:         c.Prefix,
		MiddleName:     c.MiddleName,
		Suffix:         c.Suffix,
		Birthday:       birthday,
		Phone:          c.Phone,
		Email:          c.Email,
		AdditionalInfo: additionalInfo,
		UserId:         userId,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      c.UpdatedAt.Format(time.RFC3339),
		DeletedAt:      c.DeletedAt.Time.Format(time.RFC3339),
	}
}
