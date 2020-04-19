DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    email       varchar(50) unique not null,
    nickname    varchar(50) unique not null primary key,
    fullname    varchar(60) not null,
    about       text
);

CREATE TABLE forums (
    slug    varchar(50) unique not null primary key,
    admin   varchar(50) not null,
    title   varchar(50) not null,
    FOREIGN KEY (admin) REFERENCES "users" (nickname)
);

CREATE TABLE threads (
    id      serial not null primary key,
    author  varchar(50) not null,
    created timestamp not null,
    forum   varchar(50) not null,
    message text not null,
    slug    varchar(50) unique,
    title   varchar(50) not null,
    FOREIGN KEY (forum)     REFERENCES  "forums"    (slug),
    FOREIGN KEY (author)    REFERENCES  "users"     (nickname)
);

CREATE TABLE posts (
    id          serial not null primary key,
    author      varchar(50) not null,
    forum       varchar(50) not null,
    created     timestamp not null,
    message     text not null,
    isEdited    boolean default false,
    path        integer[] not null,
    thread      integer not null,
    FOREIGN KEY (author)   REFERENCES  "users"      (nickname),
    FOREIGN KEY (thread)   REFERENCES  "threads"    (id),
    FOREIGN KEY (forum)    REFERENCES  "forums"     (slug)
);

CREATE TABLE votes (
    id      serial  not null primary key,
    thread  integer not null,
    author  varchar(50)  not null,
    vote    bool    not null,
    FOREIGN KEY (thread)    REFERENCES  "threads"   (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (nickname)
);
