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

func (s *CustomerService) Create(ctx context.Context, customer *entity.Customer) error {
	return s.DB.Create(customer).Error
}

func (s *CustomerService) Get(ctx context.Context, id int64) (*entity.Customer, error) {
	var customer entity.Customer
	if err := s.DB.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) Update(ctx context.Context, customer *entity.Customer) error {
	return s.DB.Save(customer).Error
}

func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.Customer{}, "id = ?", id).Error
}

func (s *CustomerService) List(ctx context.Context, limit, offset int) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	var total int64

	s.DB.Model(&entity.Customer{}).Count(&total)
	if err := s.DB.Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
