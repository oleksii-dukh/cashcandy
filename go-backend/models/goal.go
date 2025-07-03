package models

import (
	"database/sql"
	"time"
)

type Goal struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	TargetAmount  float64   `json:"target_amount" db:"target_amount"`
	CurrentAmount float64   `json:"current_amount" db:"current_amount"`
	Deadline      time.Time `json:"deadline" db:"deadline"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type GoalRepository struct {
	db *sql.DB
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(goal *Goal) error {
	query := `
		INSERT INTO goals (user_id, title, target_amount, current_amount, deadline, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, goal.UserID, goal.Title, goal.TargetAmount, goal.CurrentAmount, goal.Deadline, time.Now())
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	goal.ID = int(id)
	return nil
}

func (r *GoalRepository) GetByUserID(userID int) ([]Goal, error) {
	query := `
		SELECT id, user_id, title, target_amount, current_amount, deadline, created_at
		FROM goals
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.Title, &goal.TargetAmount,
			&goal.CurrentAmount, &goal.Deadline, &goal.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func (r *GoalRepository) GetByID(id int) (*Goal, error) {
	goal := &Goal{}
	query := `
		SELECT id, user_id, title, target_amount, current_amount, deadline, created_at
		FROM goals
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&goal.ID, &goal.UserID, &goal.Title, &goal.TargetAmount,
		&goal.CurrentAmount, &goal.Deadline, &goal.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (r *GoalRepository) Update(goal *Goal) error {
	query := `
		UPDATE goals 
		SET title = ?, target_amount = ?, deadline = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, goal.Title, goal.TargetAmount, goal.Deadline, goal.ID)
	return err
}

func (r *GoalRepository) Delete(id int) error {
	query := `DELETE FROM goals WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *GoalRepository) UpdateCurrentAmount(id int, amount float64) error {
	query := `UPDATE goals SET current_amount = ? WHERE id = ?`
	_, err := r.db.Exec(query, amount, id)
	return err
}
