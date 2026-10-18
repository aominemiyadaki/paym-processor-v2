package service

import (
	"context"
	"paym-processor/internal/wallet/domain"
)

type Service struct {
	walletRepo  domain.WalletRepository
	paymentRepo domain.PaymentRepository
}

func NewService(walletRepo domain.WalletRepository, paymentRepo domain.PaymentRepository) *Service {
	return &Service{
		walletRepo:  walletRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *Service) ProcessPayment(ctx context.Context, paymentId, walletId string, amount int64) error {
	p, err := domain.NewPayment(paymentId, walletId, amount)
	if err != nil {
		return err
	}

	if err := s.paymentRepo.SavePayment(ctx, p); err != nil {
		return err
	}

	w, err := s.walletRepo.GetByIDWallet(ctx, walletId)
	if err != nil {
		_ = s.paymentRepo.UpdateStatusPayment(ctx, p.ID(), domain.StatusFailed)
		return err
	}

	if err := w.Withdraw(amount); err != nil {
		_ = s.paymentRepo.UpdateStatusPayment(ctx, p.ID(), domain.StatusFailed)
		return err
	}
	if err := s.walletRepo.UpdateWallet(ctx, w); err != nil {
		_ = s.paymentRepo.UpdateStatusPayment(ctx, p.ID(), domain.StatusFailed)
		return err
	}
	if err := s.paymentRepo.UpdateStatusPayment(ctx, paymentId, domain.StatusSuccess); err != nil {
		return err
	}
	return nil
}
