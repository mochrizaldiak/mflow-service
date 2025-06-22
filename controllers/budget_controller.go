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

func RegisterBudgetRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	budgetService := services.NewBudgetService(db)

	rg.GET("/", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		data, err := budgetService.GetByUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	rg.POST("/", func(c *gin.Context) {
		var input models.Budget
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.UserID = c.GetUint("user_id")
		input.Tanggal = time.Now()
		if err := budgetService.Create(&input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, input)
	})

	rg.GET("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		userID := c.GetUint("user_id")
		b, err := budgetService.GetByID(uint(id), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, b)
	})

	rg.PUT("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		userID := c.GetUint("user_id")
		var data models.Budget
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := budgetService.Update(uint(id), userID, &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data diupdate"})
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		userID := c.GetUint("user_id")
		if err := budgetService.Delete(uint(id), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data dihapus"})
	})
}
