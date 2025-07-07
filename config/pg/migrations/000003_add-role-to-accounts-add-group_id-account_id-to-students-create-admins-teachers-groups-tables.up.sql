CREATE TABLE teachers
(
    id           int generated always as identity primary key,
    fio          varchar(256),
    phone_number varchar(20),
    account_id int references accounts (id),
    is_deleted      bool default false
);

CREATE TABLE admins
(
    id           int generated always as identity primary key,
    fio          varchar(256),
    phone_number varchar(20),
    account_id int references accounts (id)
);

CREATE TABLE groups
(
    id         int generated always as identity primary key,
    name       varchar(256),
    teacher_id int references teachers (id),
    account_id int references accounts (id),
    is_deleted    bool default false
);

ALTER TABLE accounts
    add column role varchar(20);

ALTER TABLE students
    add column group_id int references groups (id),
    add column account_id int references groups (id);