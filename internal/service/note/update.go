package note

import (
	"context"
	"github.com/vdyakova/grpc/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) Update(ctx context.Context, info *model.Note) (*emptypb.Empty, error) {
	_, err := s.noteRepository.Update(ctx, info)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
