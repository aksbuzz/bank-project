-- name: CreateLoan :one
INSERT INTO loans (book_id, member_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetLoan :one
SELECT l.id, b.id AS book_id, b.title, b.author, m.id AS member_id, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
INNER JOIN books b ON l.book_id = b.id
INNER JOIN members m ON l.member_id = m.id
WHERE l.id = $1 LIMIT 1;

-- name: ListLoans :many
SELECT l.id, b.id AS book_id, b.title, b.author, m.id AS member_id, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
INNER JOIN books b ON l.book_id = b.id
INNER JOIN members m ON l.member_id = m.id
WHERE l.id > sqlc.arg('cursor')
ORDER BY l.id
LIMIT sqlc.arg('limit')
;

-- -- name: ListOverdueLoans :many
-- SELECT b.title, b.author, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
-- INNER JOIN books b ON l.book_id = b.id
-- INNER JOIN members m ON l.member_id = m.id
-- WHERE l.return_date IS NULL AND CURDATE() > l.due_date
-- ORDER BY l.id
-- LIMIT $1
-- OFFSET $2;

-- -- name: ListLoansByMember :many
-- SELECT b.title, b.author, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
-- INNER JOIN books b ON l.book_id = b.id
-- INNER JOIN members m ON l.member_id = m.id
-- WHERE member_id = $1
-- ORDER BY l.id
-- LIMIT $2
-- OFFSET $3;

-- name: UpdateLoanReturnDate :exec
UPDATE loans
SET return_date = $2 
WHERE id = $1;

-- name: UpdateLoanOverdueFee :exec
UPDATE loans
SET overdue_fee = $2
WHERE id = $1;

-- name: DeleteLoan :exec
DELETE FROM loans
WHERE id = $1;
