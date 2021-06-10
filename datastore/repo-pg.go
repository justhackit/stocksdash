package datastore

import (
	"context"

	utils "github.com/cloudlifter/go-utils/timeutils"
	"github.com/hashicorp/go-hclog"
	"gorm.io/gorm"
)

// PostgresRepository has the implementation of the db methods.
type PostgresRepository struct {
	db     *gorm.DB
	logger hclog.Logger
}

// NewPostgresRepository returns a new PostgresRepository instance
func NewPostgresRepository(db *gorm.DB, logger hclog.Logger) *PostgresRepository {
	return &PostgresRepository{db, logger}
}

func (repo *PostgresRepository) AddNewTrade(ctx context.Context, userId string, trade *Holding) error {
	defer utils.TimeTaken("AddNewTrade", repo.logger)()
	trade.UserId = userId
	repo.logger.Info("creating user", hclog.Fmt("%#v", trade))
	result := repo.db.Debug().Create(&trade)
	return result.Error
}
