package usecase

import (
	"balance/internal/logger"
	"balance/internal/model"
	"errors"
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
		u.logger.Error("Amount must be greater than zero")
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Introduction(userID, amount)
}

func (u *Usecase) Debit(userID int, amount float64) error {
	if amount <= 0 {
		u.logger.Error("Amount must be greater than zero")
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Debit(userID, amount)
}

func (u *Usecase) Transfer(fromUserID, toUserID int, amount float64) error {
	if fromUserID == toUserID {
		u.logger.Error("Cannot transfer to the same user")
		return errors.New("cannot transfer to the same user")
	}
	if amount <= 0 {
		u.logger.Error("Amount must be greater than zero")
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Transfer(fromUserID, toUserID, amount)
}

func (u *Usecase) GetBalance(userID int) (model.Balance, error) {
	balance, err := u.repo.GetBalance(userID)
	if err != nil {
		u.logger.Error(err.Error())
		return model.Balance{}, err
	}
	u.logger.Info("Get balance completed successfully")
	return model.Balance{
		UserID:  userID,
		Balance: balance,
	}, nil
}
