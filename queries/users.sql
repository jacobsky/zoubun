-- name: SelectCounter :one
SELECT current_count
FROM counters
WHERE userid=$1;

-- name: IncrementCounter :one
UPDATE counters
SET current_count=$2
WHERE userid =$1
RETURNING current_count;
