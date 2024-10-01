package note

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) Delete(ctx context.Context, id int64) (*emptypb.Empty, error) {

	_, err := s.noteRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
