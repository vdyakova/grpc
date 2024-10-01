package note

import (
	"context"
	"log"

	"github.com/vdyakova/grpc/internal/converter"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	noteInfo := converter.ToNoteInfoFromDesc(req)
	id, err := i.noteService.Create(ctx, noteInfo)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
