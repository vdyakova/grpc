package converter

import (
	"github.com/vdyakova/grpc/internal/model"
	modelRepo "github.com/vdyakova/grpc/internal/repository/note/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(info modelRepo.NoteInfo) model.NoteInfo {
	return model.NoteInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  info.Role,
	}
}
