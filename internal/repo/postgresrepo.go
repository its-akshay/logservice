package repo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/logservice/internal/model"
	"gorm.io/gorm"
)

type postgresRepo struct {
	db *gorm.DB
}

// constructor
func NewPostgresRepo(db *gorm.DB) LogRepository {
	return &postgresRepo{db: db}
}

func (p *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Log, error) {
	var log model.Log

	err := p.db.WithContext(ctx).First(&log, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &log, nil
}
func (p *postgresRepo) List(ctx context.Context, f Filter) ([]model.Log, error) {
	var logs []model.Log

	query := p.db.WithContext(ctx).Model(&model.Log{})

	if f.Level != "" {
		query = query.Where("level = ?", f.Level)
	}

	if f.Service != "" {
		query = query.Where("service = ?", f.Service)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	if f.Offset > 0 {
		query = query.Offset(f.Offset)
	}

	err := query.Order("created_at desc").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}
func (p *postgresRepo) Insert(ctx context.Context, log *model.Log) error {
	return p.db.WithContext(ctx).Create(log).Error
}
