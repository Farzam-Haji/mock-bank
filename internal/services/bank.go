package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Farzam-Haji/mock-bank/internal/models"
)

var ErrInternal = errors.New("Internal Server Issue")

type BankRepo interface {
	CreatePayment(ctx context.Context, payment *models.Payment) (error)
}

type bankService struct {
	repo 	BankRepo
}

func NewBankService(repoInterface BankRepo) (*bankService) {
	return &bankService{
		repo: repoInterface,
	}
}

func (s *bankService) CreatePayment(ctx context.Context, orderID int64, amount int64) (*models.Payment, error) {
	payment := models.Payment{
		MerchantOrderID: orderID,
		Amount: amount,
	}

	err := s.repo.CreatePayment(ctx, &payment)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return &payment, nil
}