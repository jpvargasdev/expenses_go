package routes

import (
	"guilliman/internal/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.GET("", c.GetCategoriesController)
			categories.POST("", c.CreateCategoryController)
			// categories.PUT("/categories/:id", c.UpdateCategoryController)
			// categories.DELETE("/categories/:id", c.DeleteCategoryController)
		}
		accounts := v1.Group("/accounts")
		{
			accounts.GET("", c.GetAccountsController)
			accounts.POST("", c.AddAccountController)
			// accounts.PUT("/accounts/:id", c.UpdateAccountController)
			// accounts.DELETE("/accounts/:id", c.DeleteAccountController)
		}
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", c.GetTransactionsController)
			transactions.POST("", c.AddTransactionController)
			// transactions.PUT("/transactions/:id", c.UpdateTransactionController)
			transactions.DELETE(":id", c.DeleteTransactionController)

			// Transaccions by type
			transactions.GET("/expenses", c.GetExpensesController) // Tipo 'Expense'
			transactions.GET("/incomes", c.GetIncomesController)   // Tipo 'Income'
			transactions.GET("/savings", c.GetSavingsController)   // Tipo 'Savings'

			// Transaccions by period
			transactions.GET("/period", c.GetTransactionsForPeriodController)
			transactions.GET("/monthly", c.GetTransactionsMonthlyController)
		}
		// budget := v1.Group("/budget")
		// {
		// 	// budget.GET("/summary", c.GetBudgetSummaryController)
		// }
		transfers := v1.Group("/transfers")
		{
			transfers.GET("", c.GetTransfersController)
			transfers.POST("", c.TransferFundsController)
		}
	}
	// Health
	// r.GET("/health", c.HealthCheckController)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
