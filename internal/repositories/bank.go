package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Farzam-Haji/mock-bank/internal/models"
)

var (
	ErrRowNotFound = errors.New("Row Not Found")
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

func (r *PostgresBankRepo) GetAccount(ctx context.Context, card string) (*models.Account, error) {
	query := `
		SELECT id, name, card_number, password, balance
		FROM accounts
		WHERE card_number = $1
	`
	row := r.db.QueryRowContext(ctx, query, card)
	var account models.Account
	
	err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CardNumber,
		&account.Password,
		&account.Balance,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}