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
SELECT id, creation_date from USERS where username='$1';

-- name: GetUserById :one
SELECT username, creation_date FROM users WHERE id=$1;

-- name: CreateUser :one
WITH created_user AS (INSERT INTO users (username) VALUES ($1) RETURNING id)
INSERT INTO counters (userid) SELECT created_user.id FROM created_user
RETURNING userid;

-- name: AddUserKey :exec
INSERT INTO user_keys (userid, apikey1, apikey2) VALUES (sqlc.arg(userid)::int, sqlc.arg(apikey1)::text, sqlc.arg(apikey2)::text);

-- name: RotateUserKey1 :one
UPDATE user_keys SET apikey1=sqlc.arg(new_key)::text WHERE userid=sqlc.arg(userid)::int RETURNING apikey1;
--
-- name: RotateUserKey2 :one
UPDATE user_keys SET apikey2=sqlc.arg(new_key)::text WHERE userid=sqlc.arg(userid)::int RETURNING apikey2;

-- name: GetUserKey :one
SELECT apikey1, apikey2 FROM user_keys WHERE userid=$1;

-- name: GetUserIdFromAuth :one
SELECT userid FROM user_keys WHERE apikey1=sqlc.arg(apikey)::text OR apikey2=sqlc.arg(apikey)::text;

-- name: UsernameExists :one
SELECT CASE WHEN username=$1::text THEN TRUE ELSE FALSE END AS Exists FROM users; 
