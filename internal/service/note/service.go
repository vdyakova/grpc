package note

import (
	"github.com/vdyakova/grpc/internal/client/db"
	"github.com/vdyakova/grpc/internal/repository"
	"github.com/vdyakova/grpc/internal/repository_log"
	"github.com/vdyakova/grpc/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
	logRepository  repository_log.LogRepository
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
	logRepository repository_log.LogRepository,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
		logRepository:  logRepository,
	}
}
