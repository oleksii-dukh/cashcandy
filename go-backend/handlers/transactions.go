package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oleksii-dukh/cashcandy/go-backend/models"
)

type TransactionsHandler struct {
	transactionRepo TransactionRepository
	goalRepo        GoalRepository
}

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByGoalID(goalID int) ([]models.Transaction, error)
	GetByUserID(userID int) ([]models.Transaction, error)
	GetTotalByGoalID(goalID int) (float64, error)
}

type CreateTransactionRequest struct {
	GoalID      int     `json:"goal_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,min=0.01"`
	Description string  `json:"description"`
	Type        string  `json:"type" validate:"required,oneof=add remove"`
}

func NewTransactionsHandler(transactionRepo TransactionRepository, goalRepo GoalRepository) *TransactionsHandler {
	return &TransactionsHandler{
		transactionRepo: transactionRepo,
		goalRepo:        goalRepo,
	}
}

func (h *TransactionsHandler) CreateTransaction(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	var req CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Verify goal exists and belongs to user
	goal, err := h.goalRepo.GetByID(req.GoalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Goal not found"})
	}

	if goal.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	// Check if removing money doesn't make current amount negative
	if req.Type == "remove" && goal.CurrentAmount < req.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient funds in goal"})
	}

	transaction := &models.Transaction{
		UserID:      userID,
		GoalID:      req.GoalID,
		Amount:      req.Amount,
		Description: req.Description,
		Type:        req.Type,
	}

	if err := h.transactionRepo.Create(transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create transaction"})
	}

	// Update goal's current amount
	newAmount := goal.CurrentAmount
	if req.Type == "add" {
		newAmount += req.Amount
	} else {
		newAmount -= req.Amount
	}

	goal.CurrentAmount = newAmount
	if err := h.goalRepo.Update(goal); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update goal amount"})
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionsHandler) GetTransactionsByGoal(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	goalID, err := strconv.Atoi(c.Param("goal_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid goal ID"})
	}

	// Verify goal belongs to user
	goal, err := h.goalRepo.GetByID(goalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Goal not found"})
	}

	if goal.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	transactions, err := h.transactionRepo.GetByGoalID(goalID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionsHandler) GetUserTransactions(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	transactions, err := h.transactionRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}
