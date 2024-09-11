package main

import (
	"context"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "host.docker.internal:50051"
	noteID  = 10
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()
	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createResp, err := c.Create(ctx, &desc.CreateRequest{
		Name:            "Lera",
		Email:           "email",
		Password:        "password",
		PasswordConfirm: "passwordconfirm",
		Role:            desc.Role_USER, // 0 соответствует USER
	})
	if err != nil {
		log.Fatalf("failed to create note: %v", err)
	}

	// Вывод ID созданной заметки
	log.Printf("note created with ID: %v", createResp.Id)

	getResp, err := c.Get(ctx, &desc.GetRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Printf("Note info: %+v", getResp)
}
