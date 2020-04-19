DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id          serial  primary key,
    email       varchar(80) unique not null,
    nickname    varchar(80) unique not null,
    fullname    varchar(80) not null,
    about       text
);

CREATE TABLE forums (
    id      serial  primary key,
    slug    varchar(80) unique not null,
    admin   integer not null,
    title   varchar(80) not null,
    FOREIGN KEY (admin) REFERENCES "users" (id)
);

CREATE TABLE threads (
    id      serial not null primary key,
    author  integer not null,
    created timestamp (6) WITH TIME ZONE not null,
    forum   integer not null,
    message text not null,
    slug    varchar(80) unique,
    title   varchar(80) not null,
    FOREIGN KEY (forum)     REFERENCES  "forums"    (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (id)
);

CREATE TABLE posts (
    id          serial not null primary key,
    author      integer not null,
    forum       integer not null,
    created     timestamp (6) WITH TIME ZONE not null,
    message     text not null,
    isEdited    boolean default false,
    path        integer[] not null,
    thread      integer not null,
    FOREIGN KEY (author)   REFERENCES  "users"      (id),
    FOREIGN KEY (thread)   REFERENCES  "threads"    (id),
    FOREIGN KEY (forum)    REFERENCES  "forums"     (id)
);

CREATE TABLE votes (
    id      serial  not null primary key,
    thread  integer not null,
    author  integer  not null,
    vote    bool    not null,
    FOREIGN KEY (thread)    REFERENCES  "threads"   (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (id)
);
