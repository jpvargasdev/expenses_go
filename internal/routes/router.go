package routes

import (
	"guilliman/internal/controller"
	"guilliman/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		categories := v1.Group("/categories", middleware.AuthMiddleware())
		{
			categories.GET("", c.GetCategoriesController)
			categories.POST("", c.CreateCategoryController)
			categories.PUT("/:id", c.UpdateCategoryController)
			categories.DELETE("/:id", c.DeleteCategoryController)
		}
		accounts := v1.Group("/accounts", middleware.AuthMiddleware())
		{
			accounts.GET("", c.GetAccountsController)
			accounts.POST("", c.AddAccountController)
			accounts.PUT("/accounts/:id", c.UpdateAccountController)
			accounts.DELETE(":id", c.DeleteAccountController)
		}
		transactions := v1.Group("/transactions", middleware.AuthMiddleware())
		{
			transactions.GET("", c.GetTransactionsController)
			transactions.GET("/:id", c.GetTransactionByIdController)
			transactions.POST("", c.AddTransactionController)
			transactions.PUT("/:id", c.UpdateTransactionController)
			transactions.DELETE(":id", c.DeleteTransactionController)

			// Transaccions by type
			transactions.GET("/expenses", c.GetExpensesController)                        // Tipo 'Expense'
			transactions.GET("/incomes", c.GetIncomesController)                          // Tipo 'Income'
			transactions.GET("/savings", c.GetSavingsController)                          // Tipo 'Savings'
			transactions.GET("/category/:main_category", c.GetTransactionsByMainCategory) // Tipo 'Budget'

			// Transaccions by period
			transactions.GET("/period", c.GetTransactionsForPeriodController)
			transactions.GET("/monthly", c.GetTransactionsMonthlyController)

			// Transactions by account
			transactions.GET("/account/:id", c.GetTransactionsByAccountController)
		}
		budget := v1.Group("/budget", middleware.AuthMiddleware())
		{
			budget.GET("/summary", c.GetBudgetSummaryController)
		}
		transfers := v1.Group("/transfers", middleware.AuthMiddleware())
		{
			transfers.GET("", c.GetTransfersController)
			transfers.POST("", c.TransferFundsController)
		}
		reset := v1.Group("/reset", middleware.AuthMiddleware())
		{
			reset.POST("", c.ResetController)
		}
		user := v1.Group("/users", middleware.AuthMiddleware())
		{
			user.POST("/create", c.CreateUserController)
			// user.POST("/delete", c.DeleteUserController)
		}
	}
	// Health
	// r.GET("/health", c.HealthCheckController)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
