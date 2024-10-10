package note

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/vdyakova/grpc/internal/repository/note/converter"
	modelRepo "github.com/vdyakova/grpc/internal/repository/note/model"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vdyakova/grpc/internal/client/db"
	"github.com/vdyakova/grpc/internal/model"
	"github.com/vdyakova/grpc/internal/repository"
	"log"
	"time"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.NoteRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {

	builderInsert := squirrel.Insert("note").PlaceholderFormat(squirrel.Dollar).
		Columns("name", "email", "role", "created_at", "updated_at").
		Values(info.Name, info.Email, int64(info.Role), time.Now(), time.Now()).Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Error building SQL query: %v", err)
		return 0, err
	}

	log.Printf("Executing query: %s with args: %v", query, args)
	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		log.Printf("Error creating note: %v", err)
		return 0, err
	}

	return id, nil

}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	log.Printf("Server - Note id: %d", id)

	builderSelect := squirrel.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("note").
		Where(squirrel.Eq{"id": id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("select error: %w", err)
	}

	var note modelRepo.Note
	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&note.ID, &note.Info.Name, &note.Info.Email, &note.Info.Role, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return converter.ToNoteFromRepo(&note), nil

}

func (r *repo) Delete(ctx context.Context, id int64) (*emptypb.Empty, error) {

	builderDelete := squirrel.Delete("note_2").Where(squirrel.Eq{"id": id})
	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatal("Delete error", err)
	}
	q := db.Query{
		Name:     "note_repository.Delete",
		QueryRaw: query,
	}
	res, err := r.db.DB().ExecContext(ctx, q, args)
	if err != nil {
		log.Fatal("Delete error", err)
	}
	log.Printf("delete %d rows", res.RowsAffected())
	return &emptypb.Empty{}, nil
}

func (r *repo) Update(ctx context.Context, info *model.Note) (*emptypb.Empty, error) {
	builderUpdate := squirrel.Update("note_2").PlaceholderFormat(squirrel.Dollar).
		Set("updated_at", time.Now())

	if info.Info.Name != "" {
		builderUpdate = builderUpdate.Set("name", info.Info.Name)
	}
	if info.Info.Email != "" {
		builderUpdate = builderUpdate.Set("email", info.Info.Email)
	}

	builderUpdate = builderUpdate.Where(squirrel.Eq{"id": info.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatal("Update error", err)
		return nil, err
	}
	q := db.Query{
		Name:     "note_repository.Update",
		QueryRaw: query,
	}
	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Fatal("Update error", err)
		return nil, err
	}

	log.Printf("updated %d rows", res.RowsAffected())
	return &emptypb.Empty{}, nil
}
