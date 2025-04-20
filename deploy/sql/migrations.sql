CREATE TABLE IF NOT EXISTS "user" (
                                     id uuid primary key,
                                     email text unique not null,
                                     password text not null,
                                     salt text not null,
                                     role text not null
);


CREATE TABLE IF NOT EXISTS pvz (
                                           id uuid primary key,
                                           registration_date timestamptz not null,
                                           city text not null
);


CREATE TABLE IF NOT EXISTS reception (
                                        id uuid primary key,
                                        reception_datetime timestamptz not null,
                                        pvz_id uuid not null references pvz(id) on delete cascade,
                                        status text not null check (status in ('in_progress', 'close'))
);


CREATE TABLE IF NOT EXISTS product (
                                      id uuid primary key,
                                      received_at timestamptz not null,
                                      type text not null check (type in ('электроника', 'одежда', 'обувь')),
                                      reception_id uuid not null references reception(id) on delete cascade
);
