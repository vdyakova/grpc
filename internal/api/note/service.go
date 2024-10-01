package note

import (
	"github.com/vdyakova/grpc/internal/service"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server                     //содержит методы grpc-сервера
	noteService                    service.NoteService // реализация бизнес логики
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
