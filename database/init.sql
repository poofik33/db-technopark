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
    title   varchar(120) not null,
    FOREIGN KEY (admin) REFERENCES "users" (id)
);

CREATE TABLE threads (
    id      serial not null primary key,
    author  integer not null,
    created timestamp (6) WITH TIME ZONE not null,
    forum   integer not null,
    message text not null,
    slug    varchar(80) unique,
    title   varchar(120) not null,
    votes   integer default 0,
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
    parent      integer,
    thread      integer not null,
    FOREIGN KEY (author)   REFERENCES  "users"      (id),
    FOREIGN KEY (thread)   REFERENCES  "threads"    (id),
    FOREIGN KEY (forum)    REFERENCES  "forums"     (id)
);

CREATE TABLE votes (
    id      serial  not null primary key,
    thread  integer not null,
    author  integer  not null,
    vote    integer    not null,
    FOREIGN KEY (thread)    REFERENCES  "threads"   (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (id)
);

CREATE OR REPLACE FUNCTION add_votes_to_count() RETURNS TRIGGER AS
$add_votes_to_count$
BEGIN
    UPDATE threads
    SET votes = votes + new.vote
    WHERE id = new.thread;

    RETURN new;
END;
$add_votes_to_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_vote_in_count() RETURNS TRIGGER AS
$update_vote_in_count$
BEGIN
    UPDATE threads
    SET votes = votes - old.vote + new.vote
    WHERE id = new.thread;

    RETURN new;
END;
$update_vote_in_count$ LANGUAGE plpgsql;

CREATE TRIGGER upd_votes_count_update
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_vote_in_count();

CREATE TRIGGER upd_votes_count_insert
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE add_votes_to_count();

CREATE OR REPLACE FUNCTION add_post_path() RETURNS TRIGGER AS
$add_post_path$
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;

    RETURN new;
END;
$add_post_path$ LANGUAGE plpgsql;

CREATE TRIGGER add_post_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE add_post_path();

CREATE INDEX idx_forums_users ON forums (admin);
CREATE INDEX idx_threads_forums ON threads (forum, created);
CREATE INDEX idx_threads_users ON threads (author);
CREATE INDEX idx_posts_users ON posts (author);
CREATE INDEX idx_posts_threads_created ON posts (thread, created);
CREATE INDEX idx_posts_threads_path ON posts (thread, path);
CREATE INDEX idx_posts_threads_array ON posts (thread, (array_length(path, 1)));
CREATE INDEX idx_posts_forum ON posts (forum);
CREATE INDEX idx_votes_uesrs ON votes (author);
CREATE INDEX idx_votes_thread ON votes (thread);

CREATE INDEX idx_forums_slug ON forums (lower(slug));
CREATE INDEX idx_threads_slug ON threads (lower(slug));
CREATE INDEX idx_user_nikcname ON users (lower(nickname));
CREATE INDEX idx_posts_path_1 ON posts ((path[1]));
CREATE INDEX idx_votes_thread_username ON votes (thread, author);
