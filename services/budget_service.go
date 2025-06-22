package services

import (
	"mflow/models"

	"gorm.io/gorm"
)

type BudgetService struct {
	DB *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{DB: db}
}

func (s *BudgetService) GetByUser(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := s.DB.Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

func (s *BudgetService) GetByID(id, userID uint) (models.Budget, error) {
	var budget models.Budget
	err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error
	return budget, err
}

func (s *BudgetService) Create(b *models.Budget) error {
	return s.DB.Create(b).Error
}

func (s *BudgetService) Update(id, userID uint, data *models.Budget) error {
	var budget models.Budget
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		return err
	}
	return s.DB.Model(&budget).Updates(data).Error
}

func (s *BudgetService) Delete(id, userID uint) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Budget{}).Error
}
