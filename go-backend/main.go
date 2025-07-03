package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oleksii-dukh/cashcandy/go-backend/database"
	"github.com/oleksii-dukh/cashcandy/go-backend/handlers"
	authmiddleware "github.com/oleksii-dukh/cashcandy/go-backend/middleware"
	"github.com/oleksii-dukh/cashcandy/go-backend/models"
)

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create tables
	if err := database.CreateTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Initialize repositories
	userRepo := models.NewUserRepository(db)
	goalRepo := models.NewGoalRepository(db)
	transactionRepo := models.NewTransactionRepository(db)

	// Initialize handlers
	jwtKey := []byte("your-secret-key-change-this-in-production")
	authHandler := handlers.NewAuthHandler(userRepo, jwtKey)
	goalsHandler := handlers.NewGoalsHandler(goalRepo)
	transactionsHandler := handlers.NewTransactionsHandler(transactionRepo, goalRepo)
	statsHandler := handlers.NewStatsHandler(goalRepo, transactionRepo)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes
	e.GET("/", hello)
	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(authmiddleware.JWTMiddleware(jwtKey))
	
	// Goals routes
	protected.GET("/goals", goalsHandler.GetGoals)
	protected.POST("/goals", goalsHandler.CreateGoal)
	protected.GET("/goals/:id", goalsHandler.GetGoal)
	protected.PUT("/goals/:id", goalsHandler.UpdateGoal)
	protected.DELETE("/goals/:id", goalsHandler.DeleteGoal)
	
	// Transactions routes
	protected.POST("/transactions", transactionsHandler.CreateTransaction)
	protected.GET("/transactions", transactionsHandler.GetUserTransactions)
	protected.GET("/goals/:goal_id/transactions", transactionsHandler.GetTransactionsByGoal)
	
	// Stats routes
	protected.GET("/dashboard", statsHandler.GetDashboardStats)

	// Start server
	log.Println("Server starting on :1323")
	e.Logger.Fatal(e.Start(":1323"))
}