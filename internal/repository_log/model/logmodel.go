package model

import "time"

type LogModel struct {
	UserId    int64
	Action    string
	Log       string
	Timestamp time.Time
}
