package controllers

import (
	"fmt"
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
		fmt.Println("DEBUG - userID dari token:", userID)

		user, err := userService.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"no_hp": user.NoHP,
			"saldo": user.Saldo,
		})
	})
}
