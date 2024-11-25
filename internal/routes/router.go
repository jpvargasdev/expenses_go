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
  router.DELETE("/expenses/:id", handlers.RemoveExpenseHandler)
  router.GET("/expenses/period", handlers.GetExpensesForPeriodHandler)

  router.GET("/incomes", handlers.GetIncomesHandler)
  router.POST("add-income", handlers.AddIncomeHandler)
  router.DELETE("/incomes/:id", handlers.RemoveIncomeHandler)

  return router
}
