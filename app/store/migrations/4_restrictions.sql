alter table contracts
    add payer_id varchar default 'not specified' not null;

---- create above / drop below ----

alter table contracts
    drop column payer_id;