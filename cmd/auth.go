package main

import (
	"context"
	"fmt"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

const grpcPort = 50051

type NoteV1ServerImpl struct {
	desc.UnimplementedNoteV1Server
	notes map[int64]*desc.GetResponse // Хранение заметок в памяти
}

func (s *NoteV1ServerImpl) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id := int64(len(s.notes) + 1)
	note := &desc.GetResponse{
		Id:        id,
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
	s.notes[id] = note
	return &desc.CreateResponse{Id: id}, nil
}
func (s *NoteV1ServerImpl) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Server - Note id: %d", req.GetId())
	note, exists := s.notes[req.GetId()]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "note not found")
	}
	return &desc.GetResponse{
		Id:        note.Id,
		Name:      note.Name,
		Email:     note.Email,
		Role:      note.Role,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
func (s *NoteV1ServerImpl) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_, ex := s.notes[req.GetId()]
	if !ex {
		return nil, grpc.Errorf(codes.NotFound, "note not found")
	}
	delete(s.notes, req.GetId())
	return &emptypb.Empty{}, nil
}
func (s *NoteV1ServerImpl) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	note, ex := s.notes[req.GetId()]
	if !ex {
		log.Printf("Server - Note id: %d", req.GetId())
	}
	if req.Name != nil {
		note.Name = req.Name.Value
	}
	return &emptypb.Empty{}, nil
}
func main() {

	//lis, err := net.Listen("tcp", "127.0.0.1:50051")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := &NoteV1ServerImpl{
		notes: make(map[int64]*desc.GetResponse),
	}
	desc.RegisterNoteV1Server(grpcServer, server)
	log.Printf("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
