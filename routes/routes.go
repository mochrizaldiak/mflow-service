package routes

import (
	"mflow/controllers"
	"mflow/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userGroup := r.Group("/users")
	controllers.RegisterUserRoutes(userGroup, db)

	articleGroup := r.Group("/articles")
	controllers.RegisterArticleRoutes(articleGroup, db)

	budgetGroup := r.Group("/budgets")
	budgetGroup.Use(middleware.AuthMiddleware())
	controllers.RegisterBudgetRoutes(budgetGroup, db)

	transactionGroup := r.Group("/transactions")
	transactionGroup.Use(middleware.AuthMiddleware())
	controllers.RegisterTransactionRoutes(transactionGroup, db)
}
