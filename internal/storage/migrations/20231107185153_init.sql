-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_uid VARCHAR PRIMARY KEY,
    track_number VARCHAR,
    entry VARCHAR,
    locale VARCHAR,
    customer_id VARCHAR,
    date_created TIMESTAMP,
    oof_shard VARCHAR
);

CREATE TABLE delivery (
    order_uid VARCHAR PRIMARY KEY,
    name VARCHAR,
    phone VARCHAR,
    zip VARCHAR,
    city VARCHAR,
    address VARCHAR,
    region VARCHAR,
    email VARCHAR,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE payment (
    order_uid VARCHAR PRIMARY KEY,
    transaction VARCHAR,
    request_id VARCHAR,
    currency VARCHAR,
    provider VARCHAR,
    amount INT,
    payment_dt TIMESTAMP,
    bank VARCHAR,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE order_items (
    chrt_id INT UNIQUE PRIMARY KEY ,
    track_number VARCHAR,
    price INT,
    rid VARCHAR,
    name VARCHAR,
    sale INT,
    size VARCHAR,
    total_price INT,
    nm_id INT,
    brand VARCHAR,
    status INT,
    FOREIGN KEY (track_number) REFERENCES orders(track_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
