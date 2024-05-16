-- name: CreateBook :one
INSERT INTO books ( title, author, year, quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: GetBookForUpdate :one
SELECT * FROM books
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: ListBooks :many
SELECT * FROM books
WHERE id > sqlc.arg('cursor')
ORDER BY id
LIMIT sqlc.arg('limit')
;

-- name: UpdateBookQuantity :exec
UPDATE books
SET quantity = $2
WHERE id = $1;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;