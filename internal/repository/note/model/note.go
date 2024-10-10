package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64        `db:"id"`
	Info      NoteInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type NoteInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int    `db:"role"`
}
