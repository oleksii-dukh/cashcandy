package models

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	GoalID      int       `json:"goal_id" db:"goal_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"` // "add" or "remove"
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *Transaction) error {
	query := `
		INSERT INTO transactions (user_id, goal_id, amount, description, type, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, transaction.UserID, transaction.GoalID, transaction.Amount, transaction.Description, transaction.Type, time.Now())
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	transaction.ID = int(id)
	return nil
}

func (r *TransactionRepository) GetByGoalID(goalID int) ([]Transaction, error) {
	query := `
		SELECT id, user_id, goal_id, amount, description, type, created_at
		FROM transactions
		WHERE goal_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID, &transaction.UserID, &transaction.GoalID,
			&transaction.Amount, &transaction.Description, &transaction.Type,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByUserID(userID int) ([]Transaction, error) {
	query := `
		SELECT id, user_id, goal_id, amount, description, type, created_at
		FROM transactions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID, &transaction.UserID, &transaction.GoalID,
			&transaction.Amount, &transaction.Description, &transaction.Type,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetTotalByGoalID(goalID int) (float64, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'add' THEN amount ELSE -amount END), 0) as total
		FROM transactions
		WHERE goal_id = ?
	`
	var total float64
	err := r.db.QueryRow(query, goalID).Scan(&total)
	return total, err
}
