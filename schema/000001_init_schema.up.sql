CREATE TABLE balances
(
    balance_id serial NOT NULL UNIQUE,
    amount BIGINT DEFAULT 0
);

CREATE TABLE freezes
(
    freeze_id serial NOT NULL UNIQUE,
    balance_id INT REFERENCES balances (balance_id) ON DELETE CASCADE NOT NULL,
    freezed_amount BIGINT NOT NULL,
    is_approved BOOLEAN DEFAULT NULL

);
