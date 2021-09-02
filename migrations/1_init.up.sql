CREATE TABLE IF NOT EXISTS contracts
(
    id            integer default nextval('contracts_id_seq1'::regclass) not null
        constraint contracts_pk
            primary key,
    user_id       varchar                                           not null,
    flight_number varchar(10)                                                  not null,
    date          timestamp with time zone                               not null,
    ticket_price  numeric                                                not null,
    fee           numeric                                                not null,
    create_tx     text,
    payment       boolean default false,
    flight_date   integer                                                not null
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS flights
(
    id        serial
        constraint flights_pkey
            primary key,
    flight_id varchar(10) not null,
    runAt   timestamp   not null,
    name      varchar     not null
);

create table if not exists payments
(
    id          integer default nextval('table_name_id_seq'::regclass) not null
        constraint table_name_pk
            primary key,
    contract_id integer                                                not null,
    pay_system  varchar                                                not null,
    customer_id varchar                                                not null
);