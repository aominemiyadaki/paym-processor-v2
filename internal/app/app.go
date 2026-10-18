package app

import (
	"context"
	"fmt"
	"log"

	"paym-processor/internal/config"
	"paym-processor/internal/wallet/domain"
	"paym-processor/internal/wallet/repository/postgres"
	"paym-processor/internal/wallet/service"
	pgxpkg "paym-processor/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg         *config.Config
	pool        *pgxpool.Pool
	walletRepo  domain.WalletRepository
	paymentRepo domain.PaymentRepository
	service     *service.Service
}

func New(ctx context.Context) (*App, error) {
	// 1. Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// 2. Пул соединений PostgreSQL
	pool, err := pgxpkg.NewPGXPool(ctx, cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("postgres pool error: %w", err)
	}

	// 3. Сборка слоёв (Dependency Injection)
	walletRepo := postgres.NewWalletPostgresRepo(pool)
	paymentRepo := postgres.NewPaymentPostgresRepo(pool)
	svc := service.NewService(walletRepo, paymentRepo)

	return &App{
		cfg:         cfg,
		pool:        pool,
		walletRepo:  walletRepo,
		paymentRepo: paymentRepo,
		service:     svc,
	}, nil
}

// Run запускает бизнес-сценарий (позже здесь будет запуск HTTP-сервера)
func (a *App) Run(ctx context.Context) error {
	log.Println("Application initialized successfully")

	walletID := uuid.New().String()
	paymentID := uuid.New().String()

	// Тестовый сценарий создания и списания
	wallet, err := domain.NewWallet(walletID, "RUB", 1000)
	if err != nil {
		return fmt.Errorf("failed to create domain wallet: %w", err)
	}

	if err := a.walletRepo.SaveWallet(ctx, wallet); err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}
	log.Printf("Wallet created [%s], balance: %d RUB\n", walletID, wallet.Balance())

	if err := a.service.ProcessPayment(ctx, paymentID, walletID, 250); err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}
	log.Printf("Payment processed [%s] for 250 RUB\n", paymentID)

	updatedWallet, err := a.walletRepo.GetByIDWallet(ctx, walletID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated wallet: %w", err)
	}
	log.Printf("Current DB balance: %d RUB\n", updatedWallet.Balance())

	return nil
}

// Close освобождает ресурсы (Graceful Shutdown)
func (a *App) Close() {
	if a.pool != nil {
		a.pool.Close()
		log.Println("Postgres connection pool closed")
	}
}
