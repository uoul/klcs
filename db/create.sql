--------------------------------------------------------------------------
-- Schema
--------------------------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS klcs AUTHORIZATION CURRENT_USER;

--------------------------------------------------------------------------
-- Entities
--------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS klcs."user" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(100) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS klcs.role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) NOT NULL UNIQUE,
    CHECK (
        name = 'ADMIN'
        OR name = 'SELLER'
    )
);

CREATE TABLE IF NOT EXISTS klcs."shop" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS klcs.account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    locked BOOLEAN DEFAULT false,
    holder_name VARCHAR(100) DEFAULT '',
    external_id UUID
);

CREATE TABLE IF NOT EXISTS klcs.article (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(200) DEFAULT '',
    price INTEGER NOT NULL, -- price in cents
    category VARCHAR(30),
    stock_amount INTEGER,
    shop_id UUID NOT NULL,
    printer_id UUID
);

CREATE TABLE IF NOT EXISTS klcs.printer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    shop_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS klcs.transaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(10) NOT NULL,
    amount INTEGER NOT NULL,
    description VARCHAR(100),
    account_id UUID,
    user_id UUID NOT NULL,
    CHECK (type = 'CASH' OR type = 'CARD'),
    CHECK (type = 'CASH' AND account_id IS NULL OR type = 'CARD' AND account_id IS NOT NULL)
);

--------------------------------------------------------------------------
-- Relationship tables
--------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS klcs.user_shop_role (
    user_id UUID,
    shop_id UUID,
    role_id UUID,
    PRIMARY KEY (user_id, shop_id, role_id)
);

CREATE TABLE IF NOT EXISTS klcs.article_transaction (
    article_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    pieces INTEGER DEFAULT 1,
    printer_ack BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (article_id, transaction_id)
);

--------------------------------------------------------------------------
-- Relationships
--------------------------------------------------------------------------
ALTER TABLE klcs.user_shop_role ADD CONSTRAINT fk_user_shop_role_user FOREIGN KEY (user_id) REFERENCES klcs."user" (id) ON DELETE CASCADE;
ALTER TABLE klcs.user_shop_role ADD CONSTRAINT fk_user_shop_role_shop FOREIGN KEY (shop_id) REFERENCES klcs."shop" (id) ON DELETE CASCADE;
ALTER TABLE klcs.user_shop_role ADD CONSTRAINT fk_user_shop_role_role FOREIGN KEY (role_id) REFERENCES klcs.role (id) ON DELETE CASCADE;
ALTER TABLE klcs.article ADD CONSTRAINT fk_article_shop FOREIGN KEY (shop_id) REFERENCES klcs."shop" (id) ON DELETE CASCADE;
ALTER TABLE klcs.article ADD CONSTRAINT fk_article_printer FOREIGN KEY (printer_id) REFERENCES klcs.printer (id) ON DELETE SET NULL;
ALTER TABLE klcs.transaction ADD CONSTRAINT fk_transaction_account FOREIGN KEY (account_id) REFERENCES klcs.account (id) ON DELETE RESTRICT;
ALTER TABLE klcs.transaction ADD CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES klcs."user" (id) ON DELETE RESTRICT;
ALTER TABLE klcs.article_transaction ADD CONSTRAINT fk_article_transaction_article FOREIGN KEY (article_id) REFERENCES klcs.article (id) ON DELETE RESTRICT;
ALTER TABLE klcs.article_transaction ADD CONSTRAINT fk_article_transaction_transaction FOREIGN KEY (transaction_id) REFERENCES klcs.transaction (id) ON DELETE CASCADE;
ALTER TABLE klcs.printer ADD CONSTRAINT fk_printer_shop FOREIGN KEY (shop_id) REFERENCES klcs.shop (id) ON DELETE CASCADE;

--------------------------------------------------------------------------
-- Index
--------------------------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username ON klcs.user (username)
