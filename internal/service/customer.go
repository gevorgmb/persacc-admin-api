package service

import (
	"context"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type CustomerService struct {
	DB *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{DB: db}
}

func (s *CustomerService) Create(ctx context.Context, customer *entity.Customer, organizationID int64) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(customer).Error; err != nil {
			return err
		}
		orgCustomer := entity.OrganizationCustomer{
			OrganizationID: organizationID,
			CustomerID:     customer.ID,
		}
		return tx.Create(&orgCustomer).Error
	})
}

func (s *CustomerService) Get(ctx context.Context, id int64, organizationID int64) (*entity.Customer, error) {
	var customer entity.Customer
	err := s.DB.Joins("JOIN organization_customers ON organization_customers.customer_id = customers.id").
		Where("customers.id = ? AND organization_customers.organization_id = ?", id, organizationID).
		First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) Update(ctx context.Context, customer *entity.Customer, organizationID int64) error {
	// Verify relationship exists
	var count int64
	s.DB.Model(&entity.OrganizationCustomer{}).
		Where("customer_id = ? AND organization_id = ?", customer.ID, organizationID).
		Count(&count)
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return s.DB.Save(customer).Error
}

func (s *CustomerService) Delete(ctx context.Context, id int64, organizationID int64) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("customer_id = ? AND organization_id = ?", id, organizationID).
			Delete(&entity.OrganizationCustomer{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Customer{}, "id = ?", id).Error
	})
}

func (s *CustomerService) List(ctx context.Context, limit, offset int, organizationID int64) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	var total int64

	query := s.DB.Model(&entity.Customer{}).
		Joins("JOIN organization_customers ON organization_customers.customer_id = customers.id").
		Where("organization_customers.organization_id = ?", organizationID)

	query.Count(&total)
	if err := query.Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
