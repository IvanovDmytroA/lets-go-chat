package repository

import (
	"context"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	repository "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

var mr messagesRepository

type messagesRepository struct {
	repository.Worker
}

func InitMessagesRepository(w *repository.Worker) {
	mr = messagesRepository{*w}
}

func GetMessagesRepository() *messagesRepository {
	return &mr
}

func (r *messagesRepository) SaveMessage(message model.Message) error {
	ctx := context.Background()
	_, err := r.Worker.Get().NewInsert().Model(&message).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *messagesRepository) GetAllMessages() ([]model.Message, error) {
	ctx := context.Background()
	messages := make([]model.Message, 0)
	if err := r.Worker.Get().NewSelect().Model(&messages).Scan(ctx); err != nil {
		return messages, err
	}
	return messages, nil
}
