package services

import (
	"mflow/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	err := s.DB.Find(&users).Error
	return users, err
}

func (s *UserService) GetByID(id uint) (models.User, error) {
	var user models.User
	err := s.DB.First(&user, id).Error
	return user, err
}

func (s *UserService) Create(user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if user.Saldo == 0 {
		user.Saldo = 0
	}

	return s.DB.Create(user).Error
}

func (s *UserService) Delete(id uint) error {
	return s.DB.Delete(&models.User{}, id).Error
}

func (s *UserService) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	return user, err
}
