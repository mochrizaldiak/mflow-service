package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nama     string `json:"nama,omitempty"`
	Email    string `gorm:"unique" json:"email"`
	NoHP     string `json:"no_hp,omitempty"`
	Password string `json:"password"`
	Saldo    int    `json:"saldo"`
}
