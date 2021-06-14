package datastore

import (
	"context"

	utils "github.com/cloudlifter/go-utils/timeutils"
	"github.com/hashicorp/go-hclog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	repo.logger.Info("Adding new trade", hclog.Fmt("%#v", trade))
	result := repo.db.Debug().Create(&trade)
	return result.Error
}

func (repo *PostgresRepository) GetHoldingsByUser(ctx context.Context, userId string) (*[]Holding, error) {
	defer utils.TimeTaken("GetHoldingsByUser", repo.logger)()
	allTrades := []Holding{}
	result := repo.db.Debug().Where("user_id = ?", userId).Find(&allTrades)
	return &allTrades, result.Error
}

func (repo *PostgresRepository) AddHistorical(ctx context.Context, tickers *StockPrices) error {
	defer utils.TimeTaken("AddHistorical", repo.logger)()
	repo.logger.Info("Adding a historical value", hclog.Fmt("%#v", tickers))
	result := repo.db.Debug().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&tickers)
	return result.Error
}

func (repo *PostgresRepository) AddBatchHistorical(ctx context.Context, tickers *[]StockPrices) error {
	defer utils.TimeTaken("AddBatchHistorical", repo.logger)()
	result := repo.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&tickers, 100)
	return result.Error
}