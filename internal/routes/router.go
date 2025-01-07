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

	public := r.Group("/api/v1")
	user := public.Group("/user")
	{
		user.POST("/register", c.RegisterUserController)
		user.POST("/login", c.LoginUserController)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		categories := protected.Group("/categories")
		{
			categories.GET("", c.GetCategoriesController)
			categories.POST("", c.CreateCategoryController)
			// categories.PUT("/categories/:id", c.UpdateCategoryController)
			// categories.DELETE("/categories/:id", c.DeleteCategoryController)
		}
		accounts := protected.Group("/accounts")
		{
			accounts.GET("", c.GetAccountsController)
			accounts.POST("", c.AddAccountController)
			// accounts.PUT("/accounts/:id", c.UpdateAccountController)
			accounts.DELETE(":id", c.DeleteAccountController)
		}
		transactions := protected.Group("/transactions")
		{
			transactions.GET("", c.GetTransactionsController)
			transactions.POST("", c.AddTransactionController)
			// transactions.PUT("/transactions/:id", c.UpdateTransactionController)
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
		budget := protected.Group("/budget")
		{
			budget.GET("/summary", c.GetBudgetSummaryController)
		}
		transfers := protected.Group("/transfers")
		{
			transfers.GET("", c.GetTransfersController)
			transfers.POST("", c.TransferFundsController)
		}
		reset := protected.Group("/user")
		{
			reset.POST("/reset", c.ResetController)
		}
	}

	// Health
	// r.GET("/health", c.HealthCheckController)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
