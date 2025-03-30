// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: reset_users.sql

package database

import (
	"context"
)

const resetUsers = `-- name: ResetUsers :execrows
DELETE FROM users
`

func (q *Queries) ResetUsers(ctx context.Context) (int64, error) {
	result, err := q.db.ExecContext(ctx, resetUsers)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
