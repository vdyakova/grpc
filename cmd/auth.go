package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

const (
	dbDSN = "host=localhost port=54322 dbname=note user=note-user password=note-password sslmode=disable"
)
const grpcPort = 50051

type NoteV1ServerImpl struct {
	desc.UnimplementedNoteV1Server
	bd *pgxpool.Pool
}

func (s *NoteV1ServerImpl) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create request received: %v", req)

	// Создание SQL-запроса
	builderInsert := squirrel.Insert("note").PlaceholderFormat(squirrel.Dollar).
		Columns("name", "email", "role", "created_at", "updated_at").
		Values(req.Name, req.Email, int64(req.Role), time.Now(), time.Now())

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Error building SQL query: %v", err)
		return nil, err
	}

	log.Printf("Executing query: %s with args: %v", query, args)

	var id int64
	err = s.bd.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Printf("Execution error: %v", err)
		return nil, fmt.Errorf("execution error: %w", err)
	}

	log.Printf("Record created with ID: %d", id)
	return &desc.CreateResponse{Id: id}, nil
}

func (s *NoteV1ServerImpl) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Server - Note id: %d", req.GetId())

	builderSelect := squirrel.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("note_2").Where(squirrel.Eq{"id": req.GetId()})
	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("select error: %w", err)
	}

	var id int64
	var name, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.bd.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "note not found")
		}
		return nil, fmt.Errorf("failed to retrieve note: %w", err)
	}
	noteRole := desc.Role(role)

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      noteRole,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt.Time),
	}, nil
}
func (s *NoteV1ServerImpl) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	builderDelete := squirrel.Delete("note_2").Where(squirrel.Eq{"id": req.GetId()})
	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatal("Delete error", err)
	}
	res, err := s.bd.Exec(ctx, query, args...)
	if err != nil {
		log.Fatal("Delete error", err)
	}
	log.Printf("updated %d rows", res.RowsAffected())
	return &emptypb.Empty{}, nil
}
func (s *NoteV1ServerImpl) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := squirrel.Update("note_2").PlaceholderFormat(squirrel.Dollar).
		Set("updated_at", time.Now())

	if req.Name != nil {
		builderUpdate = builderUpdate.Set("name", req.Name.Value)
	}
	if req.Email != nil {
		builderUpdate = builderUpdate.Set("email", req.Email.Value)
	}

	builderUpdate = builderUpdate.Where(squirrel.Eq{"id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatal("Update error", err)
		return nil, err
	}

	res, err := s.bd.Exec(ctx, query, args...)
	if err != nil {
		log.Fatal("Update error", err)
		return nil, err
	}

	log.Printf("updated %d rows", res.RowsAffected())
	return &emptypb.Empty{}, nil
}
func main() {
	ctx := context.Background()
	con, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := &NoteV1ServerImpl{
		bd: con,
	}
	desc.RegisterNoteV1Server(grpcServer, server)
	log.Printf("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
