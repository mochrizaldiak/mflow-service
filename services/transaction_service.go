package services

import (
	"errors"
	"mflow/models"

	"gorm.io/gorm"
)

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{DB: db}
}

func (s *TransactionService) Create(userID uint, tx *models.Transaction) error {
	tx.UserID = userID

	var budget models.Budget
	err := s.DB.Where("id = ? AND user_id = ?", tx.BudgetID, userID).First(&budget).Error
	if err != nil {
		return errors.New("budget tidak ditemukan atau bukan milik user")
	}

	if err := s.DB.Create(tx).Error; err != nil {
		return err
	}

	if tx.Jenis == "pemasukan" {
		budget.Pemasukan += tx.Nominal
	} else if tx.Jenis == "pengeluaran" {
		budget.Pengeluaran += tx.Nominal
	} else {
		return errors.New("jenis tidak valid")
	}

	return s.DB.Save(&budget).Error
}

func (s *TransactionService) GetByUser(userID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := s.DB.
		Preload("Budget", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "deskripsi")
		}).
		Where("user_id = ?", userID).
		Order("tanggal desc").
		Find(&txs).Error
	return txs, err
}

func (s *TransactionService) GetByID(id uint, userID uint) (models.Transaction, error) {
	var tx models.Transaction
	err := s.DB.
		Preload("Budget", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "deskripsi")
		}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&tx).Error
	return tx, err
}

func (s *TransactionService) Update(id uint, userID uint, data *models.Transaction) error {
	var oldTx models.Transaction

	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&oldTx).Error; err != nil {
		return err
	}

	var budget models.Budget
	if err := s.DB.Where("id = ? AND user_id = ?", oldTx.BudgetID, userID).First(&budget).Error; err != nil {
		return err
	}

	tx := s.DB.Begin()

	if oldTx.Jenis == "pemasukan" {
		budget.Pemasukan -= oldTx.Nominal
	} else if oldTx.Jenis == "pengeluaran" {
		budget.Pengeluaran -= oldTx.Nominal
	}

	var newBudget *models.Budget = &budget
	if data.BudgetID != oldTx.BudgetID {
		newBudget = &models.Budget{}
		if err := tx.Where("id = ? AND user_id = ?", data.BudgetID, userID).First(newBudget).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if data.Jenis == "pemasukan" {
		newBudget.Pemasukan += data.Nominal
	} else if data.Jenis == "pengeluaran" {
		newBudget.Pengeluaran += data.Nominal
	} else {
		tx.Rollback()
		return errors.New("jenis tidak valid")
	}

	oldTx.Nominal = data.Nominal
	oldTx.Jenis = data.Jenis
	oldTx.Kategori = data.Kategori
	oldTx.Catatan = data.Catatan
	oldTx.Tanggal = data.Tanggal
	oldTx.BudgetID = data.BudgetID

	if err := tx.Save(&oldTx).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&budget).Error; err != nil {
		tx.Rollback()
		return err
	}
	if newBudget.ID != budget.ID {
		if err := tx.Save(newBudget).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *TransactionService) Delete(id uint, userID uint) error {
	var txData models.Transaction
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&txData).Error; err != nil {
		return err
	}

	var budget models.Budget
	if err := s.DB.Where("id = ? AND user_id = ?", txData.BudgetID, userID).First(&budget).Error; err != nil {
		return err
	}

	switch txData.Jenis {
	case "pemasukan":
		budget.Pemasukan -= txData.Nominal
	case "pengeluaran":
		budget.Pengeluaran -= txData.Nominal
	}

	tx := s.DB.Begin()
	if err := tx.Delete(&txData).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&budget).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
