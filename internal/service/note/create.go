package note

import (
	"context"

	"github.com/vdyakova/grpc/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64
	id, err := s.noteRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}
	return id, nil
}
