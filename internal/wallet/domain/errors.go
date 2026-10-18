package domain

import "errors"

var (
	ErrInsufficientFunds = errors.New("Недостаточно средств на балансе")
	ErrInvalidAmount     = errors.New("Некорректная сумма")
	ErrWalletNotFound    = errors.New("Кошелек не найден")
	ErrPaymentNotFound   = errors.New("Кошелек не найден")

	BalanceOverflow  = errors.New("Переполнение баланса по int64")
	CurrencyNotEmpty = errors.New("Валюта не должна быть пустой")
	WalletIDNotEmpty = errors.New("ID кошелька не должен быть пустым")

	PaymentIDNotEmpty = errors.New("ID платежа не должен быть пустым")
)
