DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id          serial not null primary key,
    email       varchar(40) unique not null,
    nickname    varchar(20) unique not null,
    fullname    varchar(60) not null,
    about       text
);