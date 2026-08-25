CREATE TABLE IF NOT EXISTS comments(
  id varchar(36) PRIMARY KEY,
  post_id varchar(36) NOT NULL,
  user_id varchar(36) NOT NULL,
  content TEXT NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);