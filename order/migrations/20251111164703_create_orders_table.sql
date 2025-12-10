-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID [] NOT NULL,
    total_price float8 NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT,
    order_status TEXT
);

-- +goose Down
DROP TABLE orders;