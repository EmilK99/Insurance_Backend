CREATE TABLE IF NOT EXISTS contracts
(
    id            serial            primary key,
    user_id       varchar                                           not null,
    flight_number varchar(10)                                                  not null,
    date          timestamp with time zone                               not null,
    ticket_price  numeric                                                not null,
    fee           numeric                                                not null,
    sc_account    varchar,
    sc_key        varchar,
    payment       boolean default false,
    flight_date   integer                                                not null,
    status varchar default 'pending' not null
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
    id          SERIAL            primary key,
    contract_id integer                                                not null
        constraint payments_fk
            references contracts,
    pay_system  varchar                                                not null,
    customer_id varchar                                                not null
);