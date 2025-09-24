-- name: Motd :one
SELECT config_value
FROM config
WHERE config_key = 'motd';

-- name: GetUserCounter :one
SELECT current_count::bigint
FROM counters
WHERE userid = sqlc.arg(userid)::int;

-- name: IncrementCounter :one
UPDATE counters
SET current_count = current_count + 1
WHERE userid = sqlc.arg(userid)::int
RETURNING current_count::BIGINT;

-- name: GetUserByName :one
SELECT id, creation_date from USERS where username='$1';

-- name: GetUserById :one
SELECT username, creation_date FROM users WHERE id=$1;

-- name: CreateUser :one
WITH created_user AS (INSERT INTO users (username) VALUES ($1) RETURNING id)
INSERT INTO counters (userid) SELECT created_user.id FROM created_user
RETURNING userid;

-- name: AddUserKey :exec
INSERT INTO user_keys (userid, apikey1, apikey2) 
VALUES (
    sqlc.arg(userid)::int, 
    digest(sqlc.arg(apikey1)::text, 'sha256'),
    digest(sqlc.arg(apikey2)::text, 'sha256')
);

-- name: RotateUserKey1 :one
UPDATE user_keys SET apikey1=digest(sqlc.arg(new_key)::text, 'sha256')
WHERE userid=sqlc.arg(userid)::int RETURNING apikey1;
--
-- name: RotateUserKey2 :one
UPDATE user_keys SET apikey2=apikey2=sqlc.arg(new_key)::text
WHERE userid=sqlc.arg(userid)::int RETURNING apikey2;

-- name: GetUserKey :one
SELECT apikey1, apikey2 FROM user_keys WHERE userid=$1;

-- name: GetUserIdFromAuth :one
SELECT userid
FROM user_keys
WHERE apikey1=digest(sqlc.arg(apikey)::text, 'sha256')
    OR apikey2=digest(sqlc.arg(apikey)::text, 'sha256');

-- name: UsernameExists :one
SELECT EXISTS (SELECT * FROM users WHERE username=$1::text);
