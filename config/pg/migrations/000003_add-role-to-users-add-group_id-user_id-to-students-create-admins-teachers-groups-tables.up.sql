CREATE TABLE teachers
(
    id           int primary key references users (id) on delete cascade,
    fio          varchar(256),
    phone_number varchar(20),
    is_deleted   bool default false
);

CREATE TABLE admins
(
    id           int primary key references users (id) on delete cascade,
    fio          varchar(256),
    phone_number varchar(20),
    is_deleted   bool default false
);

CREATE TABLE groups
(
    id         int generated always as identity primary key,
    name       varchar(256),
    teacher_id int references teachers (id),
    is_deleted bool default false
);

ALTER TABLE users
    add column role varchar(20);

ALTER TABLE students
    add column group_id int references groups (id);