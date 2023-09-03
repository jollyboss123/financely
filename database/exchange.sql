-- Enable pgcrypto for UUID support.
create extension if not exists pgcrypto;

create table if not exists exchange (
    from_currency_id uuid references currency(id) not null,
    to_currency_id uuid references currency(id) not null ,
    rate numeric(10, 5) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    primary key (from_currency_id, to_currency_id)
);
comment on table exchange is 'This table stores exchange rates between currencies.';
create index if not exists idx_exchange_from_currency_id on exchange(from_currency_id);
create index if not exists idx_exchange_to_currency_id on exchange(to_currency_id);
