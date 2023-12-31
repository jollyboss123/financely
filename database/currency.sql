-- Enable pgcrypto for UUID support.
create extension if not exists pgcrypto;

create table if not exists currency (
    id uuid default gen_random_uuid() primary key,
    code text check ( length(code) = 3 ) not null unique ,
    numeric_code text check ( length(code) <= 3 ) not null unique,
    fraction smallint not null,
    grapheme text not null,
    template text not null,
    decimal text default '.',
    thousand text default ',',
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);
comment on table currency is 'This table stores active currency codes according to the ISO 4217 standard and their formatting.';
create index if not exists idx_currency_code on currency(code);

insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('MYR', '458', 2, 'RM', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('USD', '840', 2, '$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('SGD', '702', 2, '$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('JPY', '392', 0, E'\u00a5', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('AUD', '036', 2, '$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('CAD', '124', 2, '$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('CNY', '156', 2, E'\u5143', '1 $', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('EUR', '978', 2, E'\u20ac', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('GBP', '826', 2, E'\u00a3', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('HKD', '344', 2, '$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('IDR', '360', 2, 'Rp', '$1', ',', '.');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('KRW', '410', 0, E'\u20a9', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('TWD', '901', 2, 'NT$', '$1', '.', ',');
insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ('VND', '704', 0, E'\u20ab', '1 $', '.', ',');
