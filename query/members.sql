-- name: CreateMember :one
INSERT INTO members (name, email, phone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMember :one
SELECT * FROM members
WHERE id = $1 LIMIT 1;

-- name: ListMembers :many
SELECT * FROM members
WHERE id > sqlc.arg('cursor')
ORDER BY id
LIMIT sqlc.arg('limit');

-- name: UpdateMember :exec
UPDATE members
SET name = $2, email = $3, phone = $4
WHERE id = $1;

-- name: UpdateMembershipType :exec
UPDATE members
SET membership_type = $2
WHERE id = $1;