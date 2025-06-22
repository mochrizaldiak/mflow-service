package controllers

import (
	"mflow/models"
	"mflow/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTransactionRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	txService := services.NewTransactionService(db)

	rg.GET("/", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		txs, err := txService.GetByUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, txs)
	})

	rg.GET("/:id", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		tx, err := txService.GetByID(uint(id), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, tx)
	})

	rg.POST("/", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var tx models.Transaction

		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if tx.Tanggal.IsZero() {
			tx.Tanggal = time.Now()
		}

		err := txService.Create(userID, &tx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tx)
	})

	rg.PUT("/:id", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		var input models.Transaction
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Tanggal.IsZero() {
			input.Tanggal = time.Now()
		}

		err = txService.Update(uint(id), userID, &input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil diperbarui"})
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		err = txService.Delete(uint(id), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
	})
}
