package service

import (
	"context"
	"errors"

	"persacc/internal/entity"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) Create(ctx context.Context, user *entity.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) Get(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	if err := s.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(ctx context.Context, user *entity.User) error {
	return s.DB.Save(user).Error
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.DB.Delete(&entity.User{}, "id = ?", id).Error
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	s.DB.Model(&entity.User{}).Count(&total)
	if err := s.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) Register(ctx context.Context, email, name string) (*entity.User, error) {
	// Check if user already exists
	var existingUser entity.User
	if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Find the 'user' role
	var role entity.Role
	if err := s.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("'user' role not found in the system")
		}
		return nil, err
	}

	user := entity.User{
		Name:   name,
		Email:  email,
		RoleID: role.ID,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
