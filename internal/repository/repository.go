package repository

import (
	"database/sql"
	"errors"
)

type Repository interface {
	Introduction(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	GetBalance(userID int) (float64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Introduction(userID int, amount float64) error {
	_, err := r.db.Exec(`INSERT INTO balance_service (user_id, balance) VALUES (?, ?)
                         ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance), updated_at = CURRENT_TIMESTAMP`, userID, amount)
	return err
}

func (r *repository) Debit(userID int, amount float64) error {
	result, err := r.db.Exec(`UPDATE balance_service SET balance = balance - ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND balance >= ?`, amount, userID, amount)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}

func (r *repository) Transfer(fromUserID, toUserID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE balance_service SET balance = balance - ? WHERE user_id = ? AND balance >= ?`, amount, fromUserID, amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO balance_service (user_id, balance) VALUES (?, ?)
                   ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance), updated_at = CURRENT_TIMESTAMP`, toUserID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetBalance(userID int) (float64, error) {
	var balance float64
	err := r.db.QueryRow(`SELECT balance FROM balance_service WHERE user_id = ?`, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found")
	}
	return balance, err
}
