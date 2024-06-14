package repository

import (
	"balance/internal/logger"
	"database/sql"
	"errors"
	"fmt"
)

type repository struct {
	db     Connection
	logger *logger.Logger
}

type Connection interface {
	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

func NewRepository(db Connection, log *logger.Logger) *repository {
	return &repository{db: db, logger: log}
}

const (
	introductionQuery = `INSERT INTO balance_service (user_id, balance) VALUES (?, ?)
	ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance), updated_at = CURRENT_TIMESTAMP`

	debitQuery = `UPDATE balance_service SET balance = balance - ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND balance >= ?`

	transferQuery  = `UPDATE balance_service SET balance = balance - ? WHERE user_id = ? AND balance >= ?`
	transferQuery2 = `INSERT INTO balance_service (user_id, balance) VALUES (?, ?)
                   ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance), updated_at = CURRENT_TIMESTAMP`

	getBalanceQuery = `SELECT balance FROM balance_service WHERE user_id = ?`
)

func (r *repository) Introduction(userID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to begin transaction for Introduction: %v", err))
		return err
	}

	defer tx.Rollback()

	if _, err = tx.Exec(introductionQuery, userID, amount); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to execute introduction query for user %d: %v", userID, err))
		return err
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to commit transaction for Introduction: %v", err))
		return err
	}

	r.logger.Info(fmt.Sprintf("Introduction completed successfully for user %d with amount %f", userID, amount))
	return nil
}

func (r *repository) Debit(userID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to begin transaction for Debit: %v", err))
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(debitQuery, amount, userID, amount)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to execute debit query for user %d: %v", userID, err))
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
		if err != nil {
			r.logger.Error(fmt.Sprintf("Failed to check rows affected for Debit: %v", err))
			return err
		}
		r.logger.Warning(fmt.Sprintf("Insufficient funds for user %d", userID))
		return errors.New("insufficient funds")
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to commit transaction for Debit: %v", err))
		return err
	}

	r.logger.Info(fmt.Sprintf("Debit completed successfully for user %d with amount %f", userID, amount))
	return nil
}

func (r *repository) Transfer(fromUserID, toUserID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to begin transaction for Transfer: %v", err))
		return err
	}

	defer tx.Rollback()

	if _, err = tx.Exec(transferQuery, amount, fromUserID, amount); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to execute transfer query for user %d: %v", fromUserID, err))
		return err
	}

	if _, err = tx.Exec(transferQuery2, toUserID, amount); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to execute transfer query for user %d: %v", toUserID, err))
		return err
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to commit transaction for Transfer: %v", err))
		return err
	}

	r.logger.Info(fmt.Sprintf("Transfer completed successfully from user %d to user %d with amount %f", fromUserID, toUserID, amount))
	return nil
}

func (r *repository) GetBalance(userID int) (float64, error) {
	var balance float64
	if err := r.db.QueryRow(getBalanceQuery, userID).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warning(fmt.Sprintf("User not found: %d", userID))
			return 0, errors.New("user not found")
		}
		r.logger.Error(fmt.Sprintf("Failed to get balance for user %d: %v", userID, err))
		return 0, err
	}

	r.logger.Info(fmt.Sprintf("Get balance completed successfully for user %d", userID))
	return balance, nil
}
