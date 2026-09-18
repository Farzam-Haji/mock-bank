ALTER TABLE payments
DROP CONSTRAINT chk_payments_status;

ALTER TABLE payments
DROP CONSTRAINT chk_payments_amount;

ALTER TABLE payments
DROP CONSTRAINT fk_payments_account_id;

DROP TABLE payments;


ALTER TABLE accounts
DROP CONSTRAINT chk_accounts_balance;

DROP TABLE accounts;