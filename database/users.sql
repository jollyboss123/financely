-- Enable pgcrypto for UUID support.
create extension if not exists pgcrypto;

create table users (
    id uuid default gen_random_uuid() primary key,
    first_name text,
    last_name text,
    email text unique,
    password text,
    verified_at timestamptz
)
