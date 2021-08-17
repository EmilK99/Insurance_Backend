CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS contracts
(
    id INTEGER NOT NULL,
    user_id         VARCHAR(42)                 NOT NULL,
    flight_number   TEXT                        NOT NULL,
    date            TIMESTAMP WITH TIME ZONE    NOT NULL,
    ticket_price    NUMERIC                     NOT NULL,
    fee             NUMERIC                     NOT NULL,
    create_tx       TEXT                        NOT NULL
                             );