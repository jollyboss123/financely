create table if not exists sessions (
    token text primary key,
    user_id uuid constraint session_user_fk references users on delete cascade,
    data bytea not null,
    expiry timestamptz not null
)
