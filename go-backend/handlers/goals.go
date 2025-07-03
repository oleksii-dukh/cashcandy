package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oleksii-dukh/cashcandy/go-backend/models"
)

type GoalsHandler struct {
	goalRepo GoalRepository
}

type GoalRepository interface {
	Create(goal *models.Goal) error
	GetByUserID(userID int) ([]models.Goal, error)
	GetByID(id int) (*models.Goal, error)
	Update(goal *models.Goal) error
	Delete(id int) error
}

type CreateGoalRequest struct {
	Title        string    `json:"title" validate:"required"`
	TargetAmount float64   `json:"target_amount" validate:"required,min=0.01"`
	Deadline     time.Time `json:"deadline" validate:"required"`
}

type UpdateGoalRequest struct {
	Title        string    `json:"title"`
	TargetAmount float64   `json:"target_amount" validate:"min=0.01"`
	Deadline     time.Time `json:"deadline"`
}

func NewGoalsHandler(goalRepo GoalRepository) *GoalsHandler {
	return &GoalsHandler{
		goalRepo: goalRepo,
	}
}

func (h *GoalsHandler) CreateGoal(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	var req CreateGoalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	goal := &models.Goal{
		UserID:        userID,
		Title:         req.Title,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0,
		Deadline:      req.Deadline,
	}

	if err := h.goalRepo.Create(goal); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create goal"})
	}

	return c.JSON(http.StatusCreated, goal)
}

func (h *GoalsHandler) GetGoals(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	goals, err := h.goalRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get goals"})
	}

	return c.JSON(http.StatusOK, goals)
}

func (h *GoalsHandler) GetGoal(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid goal ID"})
	}

	goal, err := h.goalRepo.GetByID(goalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Goal not found"})
	}

	// Check if goal belongs to the user
	if goal.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	return c.JSON(http.StatusOK, goal)
}

func (h *GoalsHandler) UpdateGoal(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid goal ID"})
	}

	goal, err := h.goalRepo.GetByID(goalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Goal not found"})
	}

	// Check if goal belongs to the user
	if goal.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	var req UpdateGoalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update only provided fields
	if req.Title != "" {
		goal.Title = req.Title
	}
	if req.TargetAmount > 0 {
		goal.TargetAmount = req.TargetAmount
	}
	if !req.Deadline.IsZero() {
		goal.Deadline = req.Deadline
	}

	if err := h.goalRepo.Update(goal); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update goal"})
	}

	return c.JSON(http.StatusOK, goal)
}

func (h *GoalsHandler) DeleteGoal(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid goal ID"})
	}

	goal, err := h.goalRepo.GetByID(goalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Goal not found"})
	}

	// Check if goal belongs to the user
	if goal.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	if err := h.goalRepo.Delete(goalID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete goal"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Goal deleted successfully"})
}
