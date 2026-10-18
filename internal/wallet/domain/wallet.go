package domain

import (
	"math"
	"sync"
)

type Wallet struct {
	id       string
	balance  int64
	currency string
	sync.Mutex
}

func NewWallet(id string, currency string, balance int64) (*Wallet, error) {
	if id == "" {
		return nil, WalletIDNotEmpty
	}
	if currency == "" {
		return nil, CurrencyNotEmpty
	}

	return &Wallet{
		id:       id,
		currency: currency,
		balance:  balance,
	}, nil
}

func (w *Wallet) ID() string { return w.id }
func (w *Wallet) Balance() int64 {
	w.Lock()
	defer w.Unlock()
	return w.balance
}
func (w *Wallet) Currency() string { return w.currency }

func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w.Lock()
	defer w.Unlock()

	if math.MaxInt64-amount < w.balance {
		return BalanceOverflow
	}

	w.balance += amount

	return nil
}

func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w.Lock()
	defer w.Unlock()

	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount

	return nil
}
