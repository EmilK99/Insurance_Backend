alter table contracts
    add type varchar not null;

---- create above / drop below ----

alter table contracts
drop column type;