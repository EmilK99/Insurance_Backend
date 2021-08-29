CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS contracts
(
    id INTEGER NOT NULL,
    user_id         VARCHAR(42)                 NOT NULL,
    flight_number   TEXT                        NOT NULL,
    date            TIMESTAMP WITH TIME ZONE    NOT NULL,
    ticket_price    NUMERIC                     NOT NULL,
    fee             NUMERIC                     NOT NULL,
    create_tx       TEXT                        NOT NULL,
    flight_date     TIMESTAMP WITH TIME ZONE    NOT NULL
                             );

CREATE TABLE IF NOT EXISTS flights
(
    id INTEGER NOT NULL,
    flight_id   VARCHAR(50)                 NOT NULL,
    runAt            TIMESTAMP WITH TIME ZONE    NOT NULL,
    name VARCHAR    NOT NULL
                                  );