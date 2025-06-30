CREATE TABLE accounts
(
    id       int generated always as identity primary key,
    login    varchar(256),
    password varchar(256),
    salt     varchar(256)
);