CREATE TABLE users
(
    id         BIGINT PRIMARY KEY,
    email      VARCHAR(100) UNIQUE NOT NULL,
    username   VARCHAR(50) UNIQUE  NOT NULL,
    password   VARCHAR(255)       NOT NULL,
    is_admin   BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);