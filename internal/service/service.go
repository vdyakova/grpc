package service

import (
	"context"
	"github.com/vdyakova/grpc/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
	Update(ctx context.Context, info *model.Note) (*emptypb.Empty, error)
	Delete(ctx context.Context, id int64) (*emptypb.Empty, error)
}
