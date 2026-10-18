CREATE SCHEMA IF NOT EXISTS paym_processor;

-- currencies
CREATE TABLE IF NOT EXISTS paym_processor.currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    iso_code VARCHAR(3) NOT NULL UNIQUE,
    version BIGINT DEFAULT 1
);

-- statuses
CREATE TABLE IF NOT EXISTS paym_processor.statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(16) NOT NULL UNIQUE
);

-- wallets
CREATE TABLE IF NOT EXISTS paym_processor.wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0),
    currency_id UUID NOT NULL,
    version BIGINT DEFAULT 1
);

-- payments
CREATE TABLE IF NOT EXISTS paym_processor.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK(amount > 0),
    status_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version BIGINT DEFAULT 1
);

-- transactions (idempotency key & ledger entries)
CREATE TABLE IF NOT EXISTS paym_processor.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    payment_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);