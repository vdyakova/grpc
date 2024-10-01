package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	noteApi "github.com/vdyakova/grpc/internal/api/note"
	noteRepository "github.com/vdyakova/grpc/internal/repository/note"
	noteService "github.com/vdyakova/grpc/internal/service/note"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	dbDSN = "host=localhost port=54322 dbname=note user=note-user password=note-password sslmode=disable"
)
const grpcPort = 50051

type NoteV1ServerImpl struct {
	desc.UnimplementedNoteV1Server
	bd *pgxpool.Pool
}

func main() {
	ctx := context.Background()
	poll, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	noteRepo := noteRepository.NewRepository(poll)
	noteSrv := noteService.NewService(noteRepo)
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, noteApi.NewImplementation(noteSrv))
	log.Println("Starting gRPC server on port ", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
