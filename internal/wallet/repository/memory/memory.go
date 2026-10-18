package memory

import (
	"context"
	"paym-processor/internal/wallet/domain"
	"sync"
)

// ==========================================
// Wallet Repository
// ==========================================

type WalletMemoryRepo struct {
	wallets map[string]*domain.Wallet
	sync.RWMutex
}

func NewWalletMemoryRepo() *WalletMemoryRepo {
	return &WalletMemoryRepo{
		wallets: make(map[string]*domain.Wallet),
	}
}

func (r *WalletMemoryRepo) GetByIDWallet(ctx context.Context, id string) (*domain.Wallet, error) {
	r.RLock()
	defer r.RUnlock()
	w, ok := r.wallets[id]
	if !ok {
		return nil, domain.ErrWalletNotFound
	}
	return w, nil
}

func (r *WalletMemoryRepo) SaveWallet(ctx context.Context, w *domain.Wallet) error {
	r.Lock()
	defer r.Unlock()
	r.wallets[w.ID()] = w
	return nil
}

func (r *WalletMemoryRepo) UpdateWallet(ctx context.Context, w *domain.Wallet) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.wallets[w.ID()]; !ok {
		return domain.ErrWalletNotFound
	}
	r.wallets[w.ID()] = w
	return nil
}

// ==========================================
// Payment Repository
// ==========================================

type PaymentMemoryRepo struct {
	payments map[string]*domain.Payment
	sync.RWMutex
}

func NewPaymentMemoryRepo() *PaymentMemoryRepo {
	return &PaymentMemoryRepo{
		payments: make(map[string]*domain.Payment),
	}
}

func (r *PaymentMemoryRepo) GetByIDPayment(ctx context.Context, id string) (*domain.Payment, error) {
	r.RLock()
	defer r.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return nil, domain.ErrPaymentNotFound
	}
	return p, nil
}

func (r *PaymentMemoryRepo) SavePayment(ctx context.Context, p *domain.Payment) error {
	r.Lock()
	defer r.Unlock()
	r.payments[p.ID()] = p
	return nil
}

func (r *PaymentMemoryRepo) UpdateStatusPayment(ctx context.Context, id string, status domain.PaymentStatus) error {
	r.Lock()
	defer r.Unlock()
	p, ok := r.payments[id]
	if !ok {
		return domain.ErrPaymentNotFound
	}

	switch status {
	case domain.StatusSuccess:
		_ = p.MarkSuccess()
	case domain.StatusFailed:
		_ = p.MarkFailed()
	}

	return nil
}
