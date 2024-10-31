package routes

import (
  "guilliman/internal/handlers"
  "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
  router := gin.Default()

  router.GET("/categories", handlers.GetCategoriesHandler)
  router.POST("/categories", handlers.CreateCategoryHandler)

  router.POST("/add-expense", handlers.AddExpenseHandler)
  router.GET("/expenses", handlers.GetExpensesHandler)
  router.GET("/expenses/period", handlers.GetExpensesForPeriodHandler)

  return router
}
