CREATE TABLE IF NOT EXISTS posts(
    id varchar(36) PRIMARY KEY,  
    title text NOT NULL,
    user_id varchar(36) NOT NULL,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)

