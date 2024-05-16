// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: loans.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createLoan = `-- name: CreateLoan :one
INSERT INTO loans (book_id, member_id)
VALUES ($1, $2)
RETURNING id, book_id, member_id, loan_date, due_date, return_date, overdue_fee
`

type CreateLoanParams struct {
	BookID   int32 `json:"book_id"`
	MemberID int32 `json:"member_id"`
}

func (q *Queries) CreateLoan(ctx context.Context, arg CreateLoanParams) (Loan, error) {
	row := q.db.QueryRow(ctx, createLoan, arg.BookID, arg.MemberID)
	var i Loan
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.MemberID,
		&i.LoanDate,
		&i.DueDate,
		&i.ReturnDate,
		&i.OverdueFee,
	)
	return i, err
}

const deleteLoan = `-- name: DeleteLoan :exec
DELETE FROM loans
WHERE id = $1
`

func (q *Queries) DeleteLoan(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteLoan, id)
	return err
}

const getLoan = `-- name: GetLoan :one
SELECT l.id, b.id AS book_id, b.title, b.author, m.id AS member_id, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
INNER JOIN books b ON l.book_id = b.id
INNER JOIN members m ON l.member_id = m.id
WHERE l.id = $1 LIMIT 1
`

type GetLoanRow struct {
	ID         int32          `json:"id"`
	BookID     int32          `json:"book_id"`
	Title      string         `json:"title"`
	Author     string         `json:"author"`
	MemberID   int32          `json:"member_id"`
	Name       string         `json:"name"`
	LoanDate   pgtype.Date    `json:"loan_date"`
	DueDate    pgtype.Date    `json:"due_date"`
	ReturnDate pgtype.Date    `json:"return_date"`
	OverdueFee pgtype.Numeric `json:"overdue_fee"`
}

func (q *Queries) GetLoan(ctx context.Context, id int32) (GetLoanRow, error) {
	row := q.db.QueryRow(ctx, getLoan, id)
	var i GetLoanRow
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.Title,
		&i.Author,
		&i.MemberID,
		&i.Name,
		&i.LoanDate,
		&i.DueDate,
		&i.ReturnDate,
		&i.OverdueFee,
	)
	return i, err
}

const listLoans = `-- name: ListLoans :many
SELECT l.id, b.id AS book_id, b.title, b.author, m.id AS member_id, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
INNER JOIN books b ON l.book_id = b.id
INNER JOIN members m ON l.member_id = m.id
WHERE l.id > $1
ORDER BY l.id
LIMIT $2
`

type ListLoansParams struct {
	Cursor int32 `json:"cursor"`
	Limit  int32 `json:"limit"`
}

type ListLoansRow struct {
	ID         int32          `json:"id"`
	BookID     int32          `json:"book_id"`
	Title      string         `json:"title"`
	Author     string         `json:"author"`
	MemberID   int32          `json:"member_id"`
	Name       string         `json:"name"`
	LoanDate   pgtype.Date    `json:"loan_date"`
	DueDate    pgtype.Date    `json:"due_date"`
	ReturnDate pgtype.Date    `json:"return_date"`
	OverdueFee pgtype.Numeric `json:"overdue_fee"`
}

func (q *Queries) ListLoans(ctx context.Context, arg ListLoansParams) ([]ListLoansRow, error) {
	rows, err := q.db.Query(ctx, listLoans, arg.Cursor, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLoansRow
	for rows.Next() {
		var i ListLoansRow
		if err := rows.Scan(
			&i.ID,
			&i.BookID,
			&i.Title,
			&i.Author,
			&i.MemberID,
			&i.Name,
			&i.LoanDate,
			&i.DueDate,
			&i.ReturnDate,
			&i.OverdueFee,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLoanOverdueFee = `-- name: UpdateLoanOverdueFee :exec
UPDATE loans
SET overdue_fee = $2
WHERE id = $1
`

type UpdateLoanOverdueFeeParams struct {
	ID         int32          `json:"id"`
	OverdueFee pgtype.Numeric `json:"overdue_fee"`
}

func (q *Queries) UpdateLoanOverdueFee(ctx context.Context, arg UpdateLoanOverdueFeeParams) error {
	_, err := q.db.Exec(ctx, updateLoanOverdueFee, arg.ID, arg.OverdueFee)
	return err
}

const updateLoanReturnDate = `-- name: UpdateLoanReturnDate :exec


UPDATE loans
SET return_date = $2 
WHERE id = $1
`

type UpdateLoanReturnDateParams struct {
	ID         int32       `json:"id"`
	ReturnDate pgtype.Date `json:"return_date"`
}

// -- name: ListOverdueLoans :many
// SELECT b.title, b.author, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
// INNER JOIN books b ON l.book_id = b.id
// INNER JOIN members m ON l.member_id = m.id
// WHERE l.return_date IS NULL AND CURDATE() > l.due_date
// ORDER BY l.id
// LIMIT $1
// OFFSET $2;
// -- name: ListLoansByMember :many
// SELECT b.title, b.author, m.name, l.loan_date, l.due_date, l.return_date, l.overdue_fee FROM loans l
// INNER JOIN books b ON l.book_id = b.id
// INNER JOIN members m ON l.member_id = m.id
// WHERE member_id = $1
// ORDER BY l.id
// LIMIT $2
// OFFSET $3;
func (q *Queries) UpdateLoanReturnDate(ctx context.Context, arg UpdateLoanReturnDateParams) error {
	_, err := q.db.Exec(ctx, updateLoanReturnDate, arg.ID, arg.ReturnDate)
	return err
}