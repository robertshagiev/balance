package repository

import (
	"balance/internal/logger"
	"database/sql"
	"errors"
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
		r.logger.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	if _, err = tx.Exec(introductionQuery, userID, amount); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	r.logger.Info("Introduction completed successfully")
	return tx.Commit()
}

func (r *repository) Debit(userID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(debitQuery, amount, userID, amount)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
		if err != nil {
			r.logger.Error(err.Error())
			return err
		}
		r.logger.Warning("Insufficient funds")
		return errors.New("insufficient funds")
	}
	r.logger.Info("Debit completed successfully")
	return tx.Commit()
}

func (r *repository) Transfer(fromUserID, toUserID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	defer tx.Rollback()

	if _, err = tx.Exec(transferQuery, amount, fromUserID, amount); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	if _, err = tx.Exec(transferQuery2, toUserID, amount); err != nil {
		r.logger.Error(err.Error())
		return err
	}

	r.logger.Info("Transfer completed successfully")
	return tx.Commit()
}

func (r *repository) GetBalance(userID int) (float64, error) {
	var balance float64
	if err := r.db.QueryRow(getBalanceQuery, userID).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warning("User not found")
			return 0, errors.New("user not found")
		}
		r.logger.Error(err.Error())
		return 0, err
	}

	r.logger.Info("Get balance completed successfully")
	return balance, nil
}
