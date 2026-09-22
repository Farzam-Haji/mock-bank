package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Farzam-Haji/mock-bank/internal/models"
	"github.com/Farzam-Haji/mock-bank/internal/repositories"
)

var (
	ErrInternal = errors.New("Internal Server Issue")
	ErrBadPayment = errors.New("Payment Cannot Proceed")
	// ErrPaymentNotFound = errors.New("Payment Not Found")
)

type BankRepo interface {
	CreatePayment(ctx context.Context, payment *models.Payment) (error)
	GetPaymentByID(ctx context.Context, paymentID int64) (*models.Payment, error)
	ProcessPayment(ctx context.Context, payment *models.Payment, account *models.Account) (error)
	UpdateFailedPayment(ctx context.Context, payment *models.Payment) (error)
}

type bankService struct {
	repo 	BankRepo
	callbackURL string
}

func NewBankService(repoInterface BankRepo, callbackURL string) (*bankService) {
	return &bankService{
		repo: repoInterface,
		callbackURL: callbackURL,
	}
}

func (s *bankService) CreatePayment(ctx context.Context, orderID int64, amount int64) (*models.Payment, error) {
	payment := models.Payment{
		MerchantOrderID: orderID,
		Amount: amount,
	}

	err := s.repo.CreatePayment(ctx, &payment)
	if err != nil {
		log.Println(fmt.Errorf("%w: %w", ErrInternal, err))
		return nil, ErrInternal
	}
	return &payment, nil
}

func (s *bankService) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)

	if err != nil {
		log.Println(fmt.Errorf("%w: %w", ErrInternal, err))
		return nil, ErrInternal
	}

	if payment.Status != "pending" {
		return payment, ErrBadPayment
	}

	return payment, nil
}

func (s *bankService) ProcessPayment(ctx context.Context, paymentID int64, card string, password string) (*models.Payment, error) {
	payment := models.Payment{
		ID: paymentID,
	}
	account := models.Account{
		CardNumber: card,
		Password: password,
	}

	errPay := s.repo.ProcessPayment(ctx, &payment, &account)
	if errPay != nil {
		errFail := s.repo.UpdateFailedPayment(ctx, &payment)
		if errFail != nil {
			log.Println(fmt.Errorf("payment %v couldnt update to 'failed' due to: %v", payment.ID, errFail))
		}
	}

	if errors.Is(errPay, repositories.ErrInsufficientBalance) ||
		errors.Is(errPay, repositories.ErrNonPending) ||
		errors.Is(errPay, repositories.ErrWrongCredential) {
			payment.Status = "failed"
			return &payment, errPay
		}

	if errPay != nil {
		payment.Status = "failed"
		return &payment, fmt.Errorf("%w: %w", ErrInternal, errPay)
	}

	return &payment, nil
}

type PaymentResult struct {
	PaymentID  		int64  		`json:"payment_id"`
	OrderID			int64 		`json:"order_id"`
	Status 			string  	`json:"status"`
}

func (s *bankService) CallbackClient(paymentID int64, orderID int64, status string) (error) {
	result := PaymentResult{
		PaymentID: paymentID,
		OrderID: orderID,
		Status: status,
	}

	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.callbackURL, bytes.NewReader(body)) 
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("callback returned status %d", resp.StatusCode)
    }

    return nil
}
