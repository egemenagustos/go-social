CREATE TABLE IF NOT EXISTS user_invitations(
    token bytea NOT NULL,
    user_id varchar(36) NOT NULL,

    PRIMARY KEY (token, user_id)
)