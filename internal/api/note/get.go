package note

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"

	desc "github.com/vdyakova/grpc/pkg/note_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := i.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, created_at: %v, updated_at: %v\n",
		noteObj.ID, noteObj.Info.Name, noteObj.Info.Email, noteObj.CreatedAt, noteObj.UpdatedAt)

	return &desc.GetResponse{
		Id:        noteObj.ID,
		Name:      noteObj.Info.Name,
		Email:     noteObj.Info.Email,
		Role:      desc.Role(noteObj.Info.Role),
		CreatedAt: timestamppb.New(noteObj.CreatedAt),
		UpdatedAt: timestamppb.New(noteObj.CreatedAt),
	}, nil
}
