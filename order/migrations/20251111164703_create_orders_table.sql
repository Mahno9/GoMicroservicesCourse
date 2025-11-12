-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID [] NOT NULL,
    total_price float8 NOT NULL,
    transaction_uuid UUID,
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