package note

import (
	"context"
	"github.com/vdyakova/grpc/internal/converter"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	noteInfo := converter.ToNoteFromUpdateRequest(req)
	_, err := i.noteService.Update(ctx, noteInfo)
	if err != nil {

		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
