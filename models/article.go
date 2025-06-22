package models

import "time"

type Article struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Judul    string    `json:"judul"`
	Tanggal  time.Time `json:"tanggal"`
	Konten   string    `json:"konten"`
	Penulis  string    `json:"penulis"`
	Kategori string    `json:"kategori"`
}
