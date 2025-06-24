package controllers

import (
	"math"
	"mflow/middleware"
	"mflow/models"
	"mflow/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userService := services.NewUserService(db)

	rg.GET("/", func(c *gin.Context) {
		users, err := userService.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	rg.GET("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		user, err := userService.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	rg.POST("/", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := userService.Create(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := userService.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})

	rg.POST("/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan password wajib diisi"})
			return
		}

		user, err := userService.FindByEmail(input.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak ditemukan"})
			return
		}

		if !services.CheckPassword(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
			return
		}

		token, err := services.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   token,
		})
	})

	authGroup := rg.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", func(c *gin.Context) {
		userID := c.GetUint("user_id")

		user, err := userService.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}

		txs, _ := services.NewTransactionService(db).GetByUser(userID)
		budgets, _ := services.NewBudgetService(db).GetByUser(userID)
		skor := calculateScore(txs, budgets)

		c.JSON(http.StatusOK, gin.H{
			"id":            user.ID,
			"nama":          user.Nama,
			"email":         user.Email,
			"no_hp":         user.NoHP,
			"saldo":         user.Saldo,
			"skor_keuangan": skor,
		})
	})
}

func calculateScore(transactions []models.Transaction, budgets []models.Budget) int {
	var pemasukan, pengeluaran int
	kategoriSet := map[string]bool{}

	for _, tx := range transactions {
		if tx.Jenis == "pemasukan" {
			pemasukan += tx.Nominal
		} else if tx.Jenis == "pengeluaran" {
			pengeluaran += tx.Nominal
		}
		kategoriSet[tx.Kategori] = true
	}

	saving := pemasukan - pengeluaran
	savingRate := 0.0
	if pemasukan > 0 {
		savingRate = float64(saving) / float64(pemasukan)
	}
	savingRate = math.Max(0, math.Min(1, savingRate))

	activeness := math.Min(float64(len(transactions))/10.0, 1.0)

	budgetScore := 0.0
	for _, b := range budgets {
		if b.TotalAnggaran == 0 {
			continue
		}
		sisa := 1.0 - float64(b.Pengeluaran)/float64(b.TotalAnggaran)
		budgetScore += math.Max(0, sisa)
	}
	if len(budgets) > 0 {
		budgetScore /= float64(len(budgets))
	}

	kategoriScore := math.Min(float64(len(kategoriSet))/5.0, 1.0)

	skor := savingRate*40 + activeness*20 + budgetScore*25 + kategoriScore*15
	return int(math.Round(skor))
}
