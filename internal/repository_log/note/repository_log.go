package note

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/vdyakova/grpc/internal/client/db"
	"log"
	"time"

	"github.com/vdyakova/grpc/internal/model"
	"github.com/vdyakova/grpc/internal/repository_log"
)

type RepositoryLog struct {
	db db.Client
}

func NewRepositoryLog(db db.Client) repository_log.LogRepository {
	return &RepositoryLog{db: db}
}

func (r *RepositoryLog) LogAction(ctx context.Context, model *model.LogModel) error {
	builderInsert := squirrel.Insert("user_logs_table").PlaceholderFormat(squirrel.Dollar).
		Columns("user_id", "log", "action", "timestamp").
		Values(model.UserId, model.Log, model.Action, time.Now()).Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {

		return err
	}
	q := db.Query{
		Name:     "user_logs_table.Create",
		QueryRaw: query,
	}
	log.Printf("log Executing query: %s with args: %v", query, args)
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {

		return fmt.Errorf("error executing SQL query: %w", err)
	}
	log.Printf("Log entry created with ID: %d", id)
	return nil
}
