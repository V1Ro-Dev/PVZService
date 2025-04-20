CREATE TABLE IF NOT EXISTS "user" (
                                     id int primary key generated always as identity,
                                     email text unique not null,
                                     password text not null,
                                     role text not null
);


CREATE TABLE IF NOT EXISTS pickup_point (
                                           id int primary key generated always as identity,
                                           registration_date date not null,
                                           city text not null
);


CREATE TABLE IF NOT EXISTS reception (
                                        id int primary key generated always as identity,
                                        reception_datetime timestamptz not null,
                                        pickup_point_id int not null references pickup_point(id) on delete cascade,
                                        status text not null check (status in ('in_progress', 'close'))
);


CREATE TABLE IF NOT EXISTS product (
                                      id int primary key generated always as identity,
                                      received_at timestamptz not null,
                                      type text not null check (type in ('электроника', 'одежда', 'обувь')),
                                      reception_id int not null references reception(id) on delete cascade
);
