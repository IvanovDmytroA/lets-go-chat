package repository

import (
	"context"
	"sync"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	repository "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

var mr messagesRepository

type messagesRepository struct {
	w  repository.Worker
	mu sync.Mutex
}

func InitMessagesRepository(w *repository.Worker) {
	mr = messagesRepository{w: *w}
}

func GetMessagesRepository() *messagesRepository {
	return &mr
}

func (r *messagesRepository) SaveMessage(message model.Message) error {
	ctx := context.Background()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.w.Get().NewInsert().Model(&message).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *messagesRepository) GetAllMessages() ([]model.Message, error) {
	ctx := context.Background()
	messages := make([]model.Message, 0)
	if err := r.w.Get().NewSelect().Model(&messages).Scan(ctx); err != nil {
		return messages, err
	}
	return messages, nil
}
