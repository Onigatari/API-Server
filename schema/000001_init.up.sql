CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    user_id int not null,
    curr_amount bigint,
    pending_amount bigint,
    last_updated timestamp,
    constraint curr_amount_non_negative check (curr_amount >= 0)
);

CREATE TABLE transactions
(
    id SERIAL PRIMARY KEY,
    user_id_from int not null REFERENCES users(id),
    user_id_to int REFERENCES users(id),
    transaction_sum bigint,
    status varchar(255) not null,
    event_type varchar(255) not null,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE service
(
    id SERIAL PRIMARY KEY,
    user_id int not null REFERENCES users(id),
    invoice bigint,
    service_id int not null,
    order_id int not null,
    status varchar(255) not null,
    created_at timestamp,
    updated_at timestamp
);