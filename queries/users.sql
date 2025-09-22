-- name: Motd :one
SELECT config_value
FROM config
WHERE config_key = 'motd';

-- name: SelectCounter :one
SELECT current_count
FROM counters
WHERE userid = $1;

-- name: IncrementCounter :one
UPDATE counters
SET current_count = current_count + 1
WHERE userid = $1
RETURNING current_count;

-- name: GetUserByName :one
SELECT id, creation_date from USERS where username=$1;
-- name: GetUserById :one
SELECT username, creation_date FROM users WHERE id=$1;

-- name: CreateUser :one
WITH created_user AS
    (INSERT INTO users (username) VALUES ($1) RETURNING id)
INSERT INTO counters (userid) SELECT created_user.id FROM created_user
RETURNING userid;
