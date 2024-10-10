package repository_log

import (
	"github.com/vdyakova/grpc/internal/model"
	"golang.org/x/net/context"
)

type LogRepository interface {
	//GetLog(ctx context.Context, id int64) ([]*model.LogModel, error)

	LogAction(ctx context.Context, log *model.LogModel) error
}
