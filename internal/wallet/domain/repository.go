package domain

import "context"

type WalletRepository interface {
	GetByIDWallet(ctx context.Context, id string) (*Wallet, error)
	SaveWallet(ctx context.Context, w *Wallet) error
	UpdateWallet(ctx context.Context, w *Wallet) error
}

type PaymentRepository interface {
	GetByIDPayment(ctx context.Context, id string) (*Payment, error)
	SavePayment(ctx context.Context, p *Payment) error
	UpdateStatusPayment(ctx context.Context, id string, status PaymentStatus) error
}
