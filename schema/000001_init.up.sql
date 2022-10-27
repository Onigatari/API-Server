CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    curr_amount BIGINT,
    pending_amount BIGINT,
    last_updated TIMESTAMP,
    CONSTRAINT curr_amount_non_negative CHECK (curr_amount >= 0)
);

CREATE TABLE transactions
(
    id SERIAL PRIMARY KEY,
    user_id_from INT NOT NULL REFERENCES users(id),
    user_id_to INT REFERENCES users(id),
    transaction_sum BIGINT,
    status VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE service
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    invoice BIGINT,
    service_id INT NOT NULL,
    order_id INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);