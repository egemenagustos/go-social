CREATE TABLE IF NOT EXISTS users(
    id varchar(36) PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    username varchar(255) NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)