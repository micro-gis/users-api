CREATE DATABASE IF NOT EXISTS users_db;
USE users_db;
create table if not exists users
(
    id           bigint auto_increment
        primary key,
    first_name   varchar(45) null,
    last_name    varchar(45) null,
    email        varchar(45) not null,
    date_created datetime    not null,
    status       varchar(45) not null,
    password     varchar(32) not null,
    constraint users_email_uindex
        unique (email)
);
INSERT INTO users_db.users (id, first_name, last_name, email, date_created, status, password) VALUES (8, 'John', 'Doe', 'email@domain.com', '2020-10-13 11:16:17', 'active', 'be6cb1069f01cd207e6484538367bd1d');