create table if not exists expenses (
    id bigserial primary key,
    title varchar(255),
    amount numeric(8,2) not null,
    transaction_date timestamp with time zone not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
--     deleted_at timestamp with time zone default current_timestamp
)
