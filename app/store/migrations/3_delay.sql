alter table contracts
    add type varchar default 'cancel' not null;

---- create above / drop below ----

alter table contracts
drop column type;