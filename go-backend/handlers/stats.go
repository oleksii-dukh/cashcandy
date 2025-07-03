package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oleksii-dukh/cashcandy/go-backend/models"
)

type StatsHandler struct {
	goalRepo        GoalRepository
	transactionRepo TransactionRepository
}

type DashboardStats struct {
	TotalSavings     float64                `json:"total_savings"`
	TotalGoals       int                    `json:"total_goals"`
	CompletedGoals   int                    `json:"completed_goals"`
	AverageProgress  float64                `json:"average_progress"`
	RecentGoals      []models.Goal          `json:"recent_goals"`
	RecentTransactions []models.Transaction `json:"recent_transactions"`
	GoalProgress     []GoalProgressStats    `json:"goal_progress"`
}

type GoalProgressStats struct {
	Goal          models.Goal `json:"goal"`
	Progress      float64     `json:"progress"`      // percentage (0-100)
	DaysRemaining int         `json:"days_remaining"`
	IsCompleted   bool        `json:"is_completed"`
}

func NewStatsHandler(goalRepo GoalRepository, transactionRepo TransactionRepository) *StatsHandler {
	return &StatsHandler{
		goalRepo:        goalRepo,
		transactionRepo: transactionRepo,
	}
}

func (h *StatsHandler) GetDashboardStats(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	// Get all user goals
	goals, err := h.goalRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get goals"})
	}

	// Get recent transactions
	transactions, err := h.transactionRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transactions"})
	}

	// Calculate statistics
	stats := h.calculateStats(goals, transactions)

	return c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) calculateStats(goals []models.Goal, transactions []models.Transaction) DashboardStats {
	stats := DashboardStats{
		RecentGoals:      make([]models.Goal, 0),
		RecentTransactions: make([]models.Transaction, 0),
		GoalProgress:     make([]GoalProgressStats, 0),
	}

	totalSavings := 0.0
	totalProgress := 0.0
	completedGoals := 0

	// Process each goal
	for _, goal := range goals {
		totalSavings += goal.CurrentAmount
		
		// Calculate progress percentage
		progress := 0.0
		if goal.TargetAmount > 0 {
			progress = (goal.CurrentAmount / goal.TargetAmount) * 100
			if progress > 100 {
				progress = 100
			}
		}
		
		totalProgress += progress
		
		// Check if goal is completed
		isCompleted := goal.CurrentAmount >= goal.TargetAmount
		if isCompleted {
			completedGoals++
		}

		// Calculate days remaining
		daysRemaining := int(goal.Deadline.Sub(getCurrentTime()).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Add to goal progress stats
		goalProgressStats := GoalProgressStats{
			Goal:          goal,
			Progress:      progress,
			DaysRemaining: daysRemaining,
			IsCompleted:   isCompleted,
		}
		stats.GoalProgress = append(stats.GoalProgress, goalProgressStats)
	}

	// Calculate average progress
	if len(goals) > 0 {
		stats.AverageProgress = totalProgress / float64(len(goals))
	}

	// Set basic stats
	stats.TotalSavings = totalSavings
	stats.TotalGoals = len(goals)
	stats.CompletedGoals = completedGoals

	// Get recent goals (last 5)
	if len(goals) > 5 {
		stats.RecentGoals = goals[:5]
	} else {
		stats.RecentGoals = goals
	}

	// Get recent transactions (last 10)
	if len(transactions) > 10 {
		stats.RecentTransactions = transactions[:10]
	} else {
		stats.RecentTransactions = transactions
	}

	return stats
}

// Helper function to get current time (can be mocked for testing)
func getCurrentTime() time.Time {
	return time.Now()
}
