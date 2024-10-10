package note

import (
	"context"
	"fmt"
	"github.com/vdyakova/grpc/internal/model"
	"log"
	"time"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Вставляем новую запись в таблицу
		var errTx error
		id, errTx = s.noteRepository.Create(ctx, info)
		if errTx != nil {
			log.Printf("Error creating note: %v", errTx)
			return errTx
		}

		logEntry := &model.LogModel{
			UserId: int(id),
			Log:    fmt.Sprintf("Note created with ID: %d", id),
			Action: "create", // Указываем действие

			Timestamp: time.Now(), // Отмечаем текущее время
		}
		err := s.logRepository.LogAction(ctx, logEntry)
		if err != nil {
			log.Printf("Error creating log: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
