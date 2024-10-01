package converter

import (
	"database/sql"
	"github.com/vdyakova/grpc/internal/model"
	desc "github.com/vdyakova/grpc/pkg/note_v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func ToNoteFromService(note *model.Note) *desc.CreateRequest {
	return &desc.CreateRequest{
		Name:  note.Info.Name,
		Email: note.Info.Email,

		Role: desc.Role(note.Info.Role), // Convert model.Role to desc.Role
	}
}

func ToNoteInfoFromDesc(info *desc.CreateRequest) *model.NoteInfo {
	return &model.NoteInfo{

		Name:  info.Name,
		Email: info.Email,
		Role:  int(info.Role), // Convert desc.Role to model.Role
	}
}
func ToNoteFromUpdateRequest(note *desc.UpdateRequest) *model.Note {
	return &model.Note{
		ID: note.Id,
		Info: model.NoteInfo{
			Name:  getValueOrDefault(note.Name),
			Email: getValueOrDefault(note.Email),
			Role:  0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{},
	}
}
func getValueOrDefault(val *wrapperspb.StringValue) string {
	if val != nil {
		return val.GetValue()
	}
	return ""
}
