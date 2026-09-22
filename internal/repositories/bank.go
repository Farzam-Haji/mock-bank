package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Farzam-Haji/mock-bank/internal/models"
)

var (
	ErrRowNotFound = errors.New("Row Not Found")
	ErrWrongCredential = errors.New("Wrong Account Information")
	ErrInsufficientBalance = errors.New("Insufficient Balance")
	ErrNonPending = errors.New("Payment Is Not Pending")
)

type PostgresBankRepo struct {
	db 		*sql.DB
}

func NewPostgresBankRepo(db *sql.DB) (*PostgresBankRepo) {
	return &PostgresBankRepo{
		db: db,
	}
}

func (r *PostgresBankRepo) CreatePayment(ctx context.Context, payment *models.Payment) (error) {
	query := `
		INSERT INTO payments(merchant_order_id, amount)
		VALUES ($1, $2)
		RETURNING id, status, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, payment.MerchantOrderID, payment.Amount)

	err := row.Scan(
		&payment.ID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}	

func (r *PostgresBankRepo) GetPaymentByID(ctx context.Context, paymentID int64) (*models.Payment, error) {
	query := `
		SELECT id, merchant_order_id, status, amount
		FROM payments
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, paymentID)
	var payment models.Payment

	if err := row.Scan(
		&payment.ID,
		&payment.MerchantOrderID,
		&payment.Status,
		&payment.Amount,
	); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PostgresBankRepo) ProcessPayment(ctx context.Context, payment *models.Payment, account *models.Account) (error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//get payment
	query := `
		SELECT merchant_order_id, status, amount
		FROM payments
		WHERE id = $1
	`
	row := tx.QueryRowContext(ctx, query, payment.ID)

	if err = row.Scan(
		&payment.MerchantOrderID,
		&payment.Status,
		&payment.Amount,
	); err != nil {
		return err
	}

	//get account
	query = `
		SELECT id, balance
		FROM accounts
		WHERE card_number = $1
			AND password = $2
	`
	row = tx.QueryRowContext(ctx, query, account.CardNumber, account.Password)

	err = row.Scan(
		&account.ID,
		&account.Balance,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrWrongCredential
	}

	if err != nil {
		return err
	}

	//check criteria
	if account.Balance < payment.Amount {
		return ErrInsufficientBalance
	}

	if payment.Status != "pending" {
		return ErrNonPending
	}

	//update account
	query = `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
			AND balance >= $1
	`
	result, err := tx.ExecContext(ctx, query, payment.Amount, account.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrInsufficientBalance
	}

	//update payment
	query = `
		UPDATE payments
		SET status = 'paid',
			updated_at = NOW(),
			account_id = $1
		WHERE id = $2
			AND status = 'pending'
	`		
	result, err = tx.ExecContext(ctx, query, account.ID, payment.ID)
	if err != nil {
		return err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNonPending
	}


	if err = tx.Commit(); err != nil {
		return err
	}

	return  nil
}

func (r *PostgresBankRepo) UpdateFailedPayment(ctx context.Context, payment *models.Payment) (error) {
	query := `
		UPDATE payments
		SET status = 'failed',
			updated_at = NOW()
		WHERE id = $1
			AND status = 'pending'
	`
	result, err := r.db.ExecContext(ctx, query, payment.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNonPending
	}

	return nil
}