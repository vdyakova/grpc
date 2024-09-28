package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64
	Info      NoteInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type NoteInfo struct {
	Name  string
	Email string
	Role  int
}

/*  id serial primary key,
    name text not null,
    email text not null,
    role int not null,
    created_at timestamp not null default now(),
    updated_at timestamp
*/
