package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/logservice/internal/model"
)

type LogRepository interface {
	Insert(ctx context.Context, log *model.Log) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Log, error)
	List(ctx context.Context, filter Filter) ([]model.Log, error)
}
