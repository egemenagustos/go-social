CREATE TABLE IF NOT EXISTS roles (
    id varchar(36) PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE, 
    level int NOT NULL DEFAULT 0,
    description text
);

INSERT INTO roles (id, name, description, level) VALUES ('01a06219-d135-73a2-bf37-74637aa83ee6', 'user', 'a user can create posts and comments', 1);

INSERT INTO roles (id, name, description, level) VALUES ('01a0621a-9d2e-700e-8dff-31b6cdbcf657', 'moderator', 'a moderator can update other users posts', 2);

INSERT INTO roles (id, name, description, level) VALUES ('01a0621b-a6ee-7646-93e3-4ededd4bfed5', 'admin', 'an admin can manage all aspects of the application', 3);