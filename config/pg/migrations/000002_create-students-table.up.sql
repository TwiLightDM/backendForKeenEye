CREATE TABLE students
(
    id           int primary key references users(id) on delete cascade,
    fio          varchar(256),
    phone_number varchar(20),
    is_deleted   bool default false
);