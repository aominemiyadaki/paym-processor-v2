package postgres

import (
	"context"
	"errors"
	"fmt"
	"paym-processor/internal/wallet/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewWalletPostgresRepo(pool *pgxpool.Pool) *WalletPostgresRepo {
	return &WalletPostgresRepo{pool: pool}
}

func (w *WalletPostgresRepo) GetByIDWallet(ctx context.Context, id string) (*domain.Wallet, error) {
	var (
		idWallet string
		balance  int64
		currency string
	)

	query := `
		SELECT w.id::text, w.balance, c.iso_code 
		FROM paym_processor.wallets w 
		JOIN paym_processor.currencies c ON w.currency_id = c.id
		WHERE w.id = $1::uuid
	`

	err := w.pool.QueryRow(ctx, query, id).Scan(&idWallet, &balance, &currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWalletNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet from db: %w", err)
	}

	wallet, err := domain.NewWallet(idWallet, currency, balance)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet entity: %w", err)
	}

	return wallet, nil
}

func (r *WalletPostgresRepo) SaveWallet(ctx context.Context, w *domain.Wallet) error {
	query := `
		INSERT INTO paym_processor.wallets (id, balance, currency_id)
		SELECT $1::uuid, $2, c.id
		FROM paym_processor.currencies c
		WHERE c.iso_code = $3
	`

	cmdTag, err := r.pool.Exec(ctx, query, w.ID(), w.Balance(), w.Currency())
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("currency '%s' not supported or not found in currencies table", w.Currency())
	}

	return nil
}

func (r *WalletPostgresRepo) UpdateWallet(ctx context.Context, w *domain.Wallet) error {
	query := `
		UPDATE paym_processor.wallets 
		SET balance = $1, updated_at = NOW() 
		WHERE id = $2::uuid
	`

	cmdTag, err := r.pool.Exec(ctx, query, w.Balance(), w.ID())
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrWalletNotFound
	}

	return nil
}
