create table if not exists payouts
(
    id          SERIAL            primary key,
    contract_id integer                                                not null
    constraint payouts_fk
    references contracts,
    pay_system  varchar                                                not null,
    customer_id varchar                                                not null,
    amount float                                                not null
);