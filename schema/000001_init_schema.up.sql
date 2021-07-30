CREATE TABLE balances
(
    balance_id serial NOT NULL UNIQUE,
    amount BIGINT DEFAULT 0 CHECK ( amount >= 0 )
);

CREATE TABLE freezes
(
    freeze_id serial NOT NULL UNIQUE,
    balance_id INT REFERENCES balances (balance_id) ON DELETE CASCADE NOT NULL,
    freezed_amount BIGINT NOT NULL CHECK ( freezed_amount >= 0 ),
    is_approved BOOLEAN DEFAULT NULL
);

INSERT INTO "balances" ("amount") VALUES (0);
INSERT INTO "balances" ("amount") VALUES (0);
INSERT INTO "balances" ("amount") VALUES (0);
INSERT INTO "balances" ("amount") VALUES (0);