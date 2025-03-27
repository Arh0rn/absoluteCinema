CREATE TABLE IF NOT EXISTS films
(
    id           UUID PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    release_year INT          NOT NULL,
    country      VARCHAR(100) NOT NULL,
    duration     INT          NOT NULL,
    budget       INT          NOT NULL,
    box_office   INT          NOT NULL

);

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email    VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT      NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL
);
