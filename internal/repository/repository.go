package repository

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vdyakova/grpc/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
	Delete(ctx context.Context, id int64) (*emptypb.Empty, error)
	Update(ctx context.Context, info *model.Note) (*emptypb.Empty, error)
}
