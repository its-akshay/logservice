package repo_test

import (
	"context"
	"testing"

	"github.com/logservice/internal/model"
	"github.com/logservice/internal/repo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&model.Log{}); err != nil {
		t.Fatal(err)
	}

	return db
}

func seedLogs(t *testing.T, r repo.LogRepository) []model.Log {
	var logs []model.Log

	for i := 0; i < 10; i++ {
		level := "INFO"
		if i >= 5 {
			level = "ERROR"
		}

		log := model.Log{
			Level:   level,
			Service: "test",
			Message: "msg",
		}

		err := r.Insert(context.Background(), &log)
		if err != nil {
			t.Fatal(err)
		}

		logs = append(logs, log)
	}

	return logs
}

func TestInsertAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewPostgresRepo(db)

	log := &model.Log{
		Level:   "INFO",
		Service: "auth",
		Message: "test",
	}

	err := r.Insert(context.Background(), log)
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.GetByID(context.Background(), log.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != log.ID {
		t.Errorf("expected %v, got %v", log.ID, got.ID)
	}
}
func TestListFilter(t *testing.T) {
	db := setupTestDB(t)
	r := repo.NewPostgresRepo(db)

	seedLogs(t, r)

	logs, err := r.List(context.Background(), repo.Filter{
		Level: "ERROR",
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(logs) != 5 {
		t.Errorf("expected 5 logs, got %d", len(logs))
	}
}