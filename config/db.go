package config

import (
	"log"
	"mflow/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// (username):(password)@tcp(port)/(dbname)?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "mflow_user:mflow_pass@tcp(127.0.0.1:3306)/mflow?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Budget{},
		&models.Transaction{},
	)

	return db
}
