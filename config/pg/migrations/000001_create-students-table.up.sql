CREATE TABLE students
(
    id           int generated always as identity primary key,
    fio          varchar(256),
    group_name   varchar(256),
    phone_number varchar(20)
);