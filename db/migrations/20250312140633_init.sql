-- migrate:up
CREATE TABLE users (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "full_name" VARCHAR(255),
  "email" VARCHAR(255),
  "password" VARCHAR(255),
  "role" VARCHAR(50),
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_deleted" BOOLEAN
);

-- migrate:down
DROP TABLE IF EXISTS users;