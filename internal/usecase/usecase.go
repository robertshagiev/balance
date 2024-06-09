package usecase

import (
	"balance/internal/model"
	"errors"
)

type Repository interface {
	Introduction(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	GetBalance(userID int) (float64, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Introduction(userID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Introduction(userID, amount)
}

func (u *Usecase) Debit(userID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Debit(userID, amount)
}

func (u *Usecase) Transfer(fromUserID, toUserID int, amount float64) error {
	if fromUserID == toUserID {
		return errors.New("cannot transfer to the same user")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return u.repo.Transfer(fromUserID, toUserID, amount)
}

func (u *Usecase) GetBalance(userID int) (model.Balance, error) {
	balance, err := u.repo.GetBalance(userID)
	if err != nil {
		return model.Balance{}, err
	}
	return model.Balance{
		UserID:  userID,
		Balance: balance,
	}, nil
}
