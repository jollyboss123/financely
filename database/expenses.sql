-- Enable pgcrypto for UUID support.
create extension if not exists pgcrypto;

create table if not exists expenses (
    id uuid default gen_random_uuid() primary key,
    title text not null,
    amount_ud bigint not null, --in user defined currency
    currency_id_ud uuid references currency(id) not null,
    currency_code_ud text,
    amount_base bigint not null, --in user base currency
    currency_id_base uuid references currency(id) not null,
    currency_code_base text,
    transaction_date timestamp with time zone not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);
comment on table expenses is 'This table stores all users expenses.';
create index if not exists idx_currency_id_loc on expenses (currency_id_ud);
create index if not exists idx_currency_id_base on expenses (currency_id_base);

