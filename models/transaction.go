package models

import "time"

type Transaction struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `json:"-"`
	BudgetID uint      `json:"budget_id"`
	Budget   Budget    `gorm:"foreignKey:BudgetID" json:"budget"`
	Nominal  int       `json:"nominal"`
	Jenis    string    `json:"jenis"`
	Kategori string    `json:"kategori"`
	Tanggal  time.Time `json:"tanggal"`
	Catatan  string    `json:"catatan,omitempty"`
}
