package services

import (
	"errors"
	"mflow/models"
	"time"

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

	// Ambil budget
	var budget models.Budget
	if err := s.DB.Where("id = ? AND user_id = ?", tx.BudgetID, userID).First(&budget).Error; err != nil {
		return errors.New("budget tidak ditemukan atau bukan milik user")
	}

	// Ambil user
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Validasi saldo cukup jika pengeluaran
	if tx.Jenis == "pengeluaran" && user.Saldo < tx.Nominal {
		return errors.New("saldo tidak mencukupi")
	}

	// Simpan transaksi
	if err := s.DB.Create(tx).Error; err != nil {
		return err
	}

	// Update budget
	if tx.Jenis == "pemasukan" {
		budget.Pemasukan += tx.Nominal
		user.Saldo += tx.Nominal
	} else if tx.Jenis == "pengeluaran" {
		budget.Pengeluaran += tx.Nominal
		user.Saldo -= tx.Nominal
	} else {
		return errors.New("jenis tidak valid")
	}

	if err := s.DB.Save(&budget).Error; err != nil {
		return err
	}
	return s.DB.Save(&user).Error
}

func (s *TransactionService) GetByUser(userID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := s.DB.
		Preload("Budget", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "deskripsi", "status", "jenis_periode", "tanggal")
		}).
		Where("user_id = ?", userID).
		Order("tanggal desc").
		Find(&txs).Error

	if err != nil {
		return nil, err
	}

	now := time.Now()

	for i := range txs {
		b := &txs[i].Budget

		if b.Status == "S" {
			continue
		}

		var endDate time.Time
		switch b.JenisPeriode {
		case "D":
			endDate = b.Tanggal.AddDate(0, 0, 1)
		case "W":
			endDate = b.Tanggal.AddDate(0, 0, 7)
		case "M":
			endDate = b.Tanggal.AddDate(0, 1, 0)
		case "Y":
			endDate = b.Tanggal.AddDate(1, 0, 0)
		default:
			continue
		}

		if now.After(endDate) {
			// Ubah status jadi selesai
			b.Status = "S"
			s.DB.Model(&models.Budget{}).Where("id = ?", b.ID).Update("status", "S")
		}
	}

	return txs, nil
}

func (s *TransactionService) GetByID(id uint, userID uint) (models.Transaction, error) {
	var tx models.Transaction
	err := s.DB.
		Preload("Budget", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama", "deskripsi")
		}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&tx).Error
	return tx, err
}

func (s *TransactionService) Update(id uint, userID uint, data *models.Transaction) error {
	var oldTx models.Transaction
	if err := s.DB.First(&oldTx, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return err
	}

	// Ambil budget lama dan baru
	var oldBudget, newBudget models.Budget
	if err := s.DB.Where("id = ? AND user_id = ?", oldTx.BudgetID, userID).First(&oldBudget).Error; err != nil {
		return err
	}
	if err := s.DB.Where("id = ? AND user_id = ?", data.BudgetID, userID).First(&newBudget).Error; err != nil {
		return err
	}

	// Ambil user
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// ===== SALDO DAN BUDGET UPDATE SECARA KONDISIONAL =====

	switch {
	case oldTx.Jenis == "pemasukan" && data.Jenis == "pengeluaran":
		// rollback pemasukan, lalu kurangi lagi untuk pengeluaran baru
		totalKurang := oldTx.Nominal + data.Nominal
		if user.Saldo < totalKurang {
			return errors.New("saldo tidak mencukupi untuk update transaksi")
		}
		user.Saldo -= totalKurang
		oldBudget.Pemasukan -= oldTx.Nominal
		newBudget.Pengeluaran += data.Nominal

	case oldTx.Jenis == "pengeluaran" && data.Jenis == "pemasukan":
		// rollback pengeluaran, lalu tambah pemasukan baru
		user.Saldo += oldTx.Nominal + data.Nominal
		oldBudget.Pengeluaran -= oldTx.Nominal
		newBudget.Pemasukan += data.Nominal

	case oldTx.Jenis == data.Jenis:
		diff := int(data.Nominal) - int(oldTx.Nominal)
		if data.Jenis == "pengeluaran" {
			if diff > 0 && user.Saldo < diff {
				return errors.New("saldo tidak mencukupi untuk update transaksi")
			}
			user.Saldo -= diff
			oldBudget.Pengeluaran -= oldTx.Nominal
			newBudget.Pengeluaran += data.Nominal
		} else if data.Jenis == "pemasukan" {
			user.Saldo += diff
			oldBudget.Pemasukan -= oldTx.Nominal
			newBudget.Pemasukan += data.Nominal
		}
	default:
		return errors.New("jenis transaksi tidak valid")
	}

	// ===== UPDATE TRANSAKSI =====
	oldTx.Nominal = data.Nominal
	oldTx.Jenis = data.Jenis
	oldTx.Kategori = data.Kategori
	oldTx.Catatan = data.Catatan
	oldTx.Tanggal = data.Tanggal
	oldTx.BudgetID = data.BudgetID

	tx := s.DB.Begin()

	if err := tx.Save(&oldTx).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&oldBudget).Error; err != nil {
		tx.Rollback()
		return err
	}
	if oldBudget.ID != newBudget.ID {
		if err := tx.Save(&newBudget).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *TransactionService) Delete(id uint, userID uint) error {
	var txData models.Transaction
	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&txData).Error; err != nil {
		return err
	}

	var budget models.Budget
	if err := s.DB.First(&budget, txData.BudgetID).Error; err != nil {
		return err
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if txData.Jenis == "pemasukan" {
		user.Saldo -= txData.Nominal
		budget.Pemasukan -= txData.Nominal
	} else if txData.Jenis == "pengeluaran" {
		user.Saldo += txData.Nominal
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
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
