-- 1. Таблица финансового журнала (Ledger)
CREATE TABLE IF NOT EXISTS paym_processor.wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    type VARCHAR(32) NOT NULL, -- 'DEBIT', 'CREDIT'
    reference_id VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс для выборки истории по кошельку
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_wallet_id_created 
ON paym_processor.wallet_transactions (wallet_id, created_at DESC);

-- 2. Таблица кэша идемпотентности
CREATE TABLE IF NOT EXISTS paym_processor.idempotency_records (
    key VARCHAR(128) PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL, -- 'PROCESSING', 'COMPLETED', 'FAILED'
    response_code INT,
    response_payload JSONB,
    locked_until TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 hours'
);

-- Индекс для TTL-очистки
CREATE INDEX IF NOT EXISTS idx_idempotency_records_expires_at 
ON paym_processor.idempotency_records (expires_at);