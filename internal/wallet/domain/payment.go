package domain

import (
	"errors"
	"time"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "PENDING"
	StatusSuccess PaymentStatus = "SUCCESS"
	StatusFailed  PaymentStatus = "FAILED"
)

var (
	ErrInvalidPaymentState = errors.New("недопустимый переход статуса платежа")
)

type Payment struct {
	id        string
	walletID  string
	amount    int64
	status    PaymentStatus
	createAt  time.Time
	updatedAt time.Time
}

func NewPayment(id string, walletID string, amount int64) (*Payment, error) {
	if id == "" {
		return nil, PaymentIDNotEmpty
	}
	if walletID == "" {
		return nil, WalletIDNotEmpty
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now().UTC()

	return &Payment{
		id:        id,
		walletID:  walletID,
		amount:    amount,
		status:    StatusPending,
		createAt:  now,
		updatedAt: now,
	}, nil
}

func (p *Payment) ID() string            { return p.id }
func (p *Payment) WalletID() string      { return p.walletID }
func (p *Payment) Amount() int64         { return p.amount }
func (p *Payment) Status() PaymentStatus { return p.status }
func (p *Payment) CreatedAt() time.Time  { return p.createAt }
func (p *Payment) UpdatedAt() time.Time  { return p.updatedAt }

func (p *Payment) MarkSuccess() error {
	if p.status != StatusPending {
		return ErrInvalidPaymentState
	}
	p.status = StatusSuccess
	p.updatedAt = time.Now().UTC()
	return nil
}

func (p *Payment) MarkFailed() error {
	if p.status != StatusPending {
		return ErrInvalidPaymentState
	}
	p.status = StatusFailed
	p.updatedAt = time.Now().UTC()
	return nil
}
