package models

import "time"

type Budget struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	Nama          string    `json:"nama"` // ✅ nama anggaran
	Pemasukan     int       `json:"pemasukan"`
	Pengeluaran   int       `json:"pengeluaran"`
	JenisAnggaran string    `json:"jenis_anggaran"`
	Deskripsi     string    `json:"deskripsi"`
	TotalAnggaran int       `json:"total_anggaran"`
	JenisPeriode  string    `json:"jenis_periode"`
	Tanggal       time.Time `json:"tanggal"`
	Status        string    `json:"status"`
}
