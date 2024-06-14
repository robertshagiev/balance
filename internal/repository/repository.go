package usecase

import (
	"balance/internal/logger"
	"balance/internal/model"
	"errors"
	"fmt"
)

type Usecase struct {
	repo   Repository
	logger *logger.Logger
}

type Repository interface {
	Introduction(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	GetBalance(userID int) (float64, error)
}

func NewUsecase(repo Repository, log *logger.Logger) *Usecase {
	return &Usecase{repo: repo, logger: log}
}

func (u *Usecase) Introduction(userID int, amount float64) error {
	if amount <= 0 {
		u.logger.Warning(fmt.Sprintf("Invalid amount for Introduction: %f", amount))
		return errors.New("amount must be greater than zero")
	}
	err := u.repo.Introduction(userID, amount)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Introduction failed for user %d: %v", userID, err))
		return err
	}
	u.logger.Info(fmt.Sprintf("Introduction completed for user %d with amount %f", userID, amount))
	return nil
}

func (u *Usecase) Debit(userID int, amount float64) error {
	if amount <= 0 {
		u.logger.Warning(fmt.Sprintf("Invalid amount for Debit: %f", amount))
		return errors.New("amount must be greater than zero")
	}
	err := u.repo.Debit(userID, amount)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Debit failed for user %d: %v", userID, err))
		return err
	}
	u.logger.Info(fmt.Sprintf("Debit completed for user %d with amount %f", userID, amount))
	return nil
}

func (u *Usecase) Transfer(fromUserID, toUserID int, amount float64) error {
	if fromUserID == toUserID {
		u.logger.Warning("Cannot transfer to the same user")
		return errors.New("cannot transfer to the same user")
	}
	if amount <= 0 {
		u.logger.Warning(fmt.Sprintf("Invalid amount for Transfer: %f", amount))
		return errors.New("amount must be greater than zero")
	}
	err := u.repo.Transfer(fromUserID, toUserID, amount)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Transfer failed from user %d to user %d: %v", fromUserID, toUserID, err))
		return err
	}
	u.logger.Info(fmt.Sprintf("Transfer completed from user %d to user %d with amount %f", fromUserID, toUserID, amount))
	return nil
}

func (u *Usecase) GetBalance(userID int) (model.Balance, error) {
	balance, err := u.repo.GetBalance(userID)
	if err != nil {
		u.logger.Error(fmt.Sprintf("GetBalance failed for user %d: %v", userID, err))
		return model.Balance{}, err
	}
	u.logger.Info(fmt.Sprintf("GetBalance completed for user %d", userID))
	return model.Balance{
		UserID:  userID,
		Balance: balance,
	}, nil
}
