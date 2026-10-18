package postgres

import (
	"context"
	"errors"
	"fmt"
	"paym-processor/internal/wallet/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPaymentPostgresRepo(pool *pgxpool.Pool) *PaymentPostgresRepo {
	return &PaymentPostgresRepo{pool: pool}
}

func (r *PaymentPostgresRepo) GetByIDPayment(ctx context.Context, id string) (*domain.Payment, error) {
	var (
		paymentID  string
		walletID   string
		amount     int64
		statusName string
	)

	query := `
		SELECT p.id::text, p.wallet_id::text, p.amount, s.name
		FROM paym_processor.payments p
		JOIN paym_processor.statuses s ON p.status_id = s.id
		WHERE p.id = $1::uuid
	`

	err := r.pool.QueryRow(ctx, query, id).Scan(&paymentID, &walletID, &amount, &statusName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment from db: %w", err)
	}

	payment, err := domain.NewPayment(paymentID, walletID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment entity: %w", err)
	}

	switch domain.PaymentStatus(statusName) {
	case domain.StatusSuccess:
		_ = payment.MarkSuccess()
	case domain.StatusFailed:
		_ = payment.MarkFailed()
	}

	return payment, nil
}

func (r *PaymentPostgresRepo) SavePayment(ctx context.Context, p *domain.Payment) error {
	query := `
		INSERT INTO paym_processor.payments (id, wallet_id, amount, status_id)
		SELECT $1::uuid, $2::uuid, $3, s.id
		FROM paym_processor.statuses s
		WHERE s.name = $4
	`

	cmd, err := r.pool.Exec(ctx, query, p.ID(), p.WalletID(), p.Amount(), string(p.Status()))
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("status '%s' not found in database", p.Status())
	}

	return nil
}

func (r *PaymentPostgresRepo) UpdateStatusPayment(ctx context.Context, id string, status domain.PaymentStatus) error {
	query := `
		UPDATE paym_processor.payments
		SET 
			status_id = (SELECT id FROM paym_processor.statuses WHERE name = $1),
			updated_at = NOW()
		WHERE id = $2::uuid
	`

	cmdTag, err := r.pool.Exec(ctx, query, string(status), id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}
