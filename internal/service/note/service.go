package note

import (
	"github.com/vdyakova/grpc/internal/repository"
	"github.com/vdyakova/grpc/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
}

func NewService(
	noteRepository repository.NoteRepository,

) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
	}
}
