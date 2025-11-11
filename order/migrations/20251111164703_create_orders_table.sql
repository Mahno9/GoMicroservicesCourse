-- +goose Up
CREATE TABLE orders (
    id serial PRIMARY KEY,
    order_uuid varchar(255) NOT NULL,
    user_uuid varchar(255) NOT NULL,
    part_uuids varchar(255) [] NOT NULL,
    total_price float8 NOT NULL,
    transaction_uuid varchar(255),
    payment_method TEXT CHECK (
        payment_method IN (
            'CARD',
            'SBP',
            'CREDIT_CARD',
            'INVESTOR_MONEY',
            'UNKNOWN'
        )
    ),
    order_status TEXT CHECK (
        order_status IN (
            'PAID',
            'CANCELLED',
            'PENDING_PAYMENT'
        )
    )
);

-- +goose Down
DROP TABLE orders;