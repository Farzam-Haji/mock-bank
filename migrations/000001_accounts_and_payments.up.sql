CREATE TABLE accounts (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name TEXT NOT NULL,
	card_number TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	balance BIGINT NOT NULL
);

ALTER TABLE accounts
ADD CONSTRAINT chk_accounts_balance
CHECK (balance >= 0);


CREATE TABLE payments (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	merchant_order_id BIGINT NOT NULL,
	account_id BIGINT,
	amount BIGINT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE payments
ADD CONSTRAINT fk_payments_account_id
FOREIGN KEY (account_id) REFERENCES accounts(id)
ON DELETE SET NULL;

ALTER TABLE payments
ADD CONSTRAINT chk_payments_amount
CHECK (amount > 0);

ALTER TABLE payments
ADD CONSTRAINT chk_payments_status
CHECK (status IN ('pending', 'paid', 'failed'));
