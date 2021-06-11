package datastore

import (
	"context"
	"testing"

	logutils "github.com/cloudlifter/go-utils/logs"
	conf "github.com/justhackit/stockedash/config"
)

var repository *PostgresRepository

func init() {
	logger := logutils.NewLogger()
	configs := conf.NewConfigurationsFromFile("../test-config.yaml", logger)
	db, err := NewConnection(configs, logger)
	if err != nil {
		logger.Error("unable to connect to db", "error", err)
		panic(err)
	}
	repository = NewPostgresRepository(db, logger)
	db.Debug().AutoMigrate(&Holding{}, &HistoricalPrices{})
	db.Exec(`ALTER TABLE HistoricalPrices
	ADD CONSTRAINT constraint_fk
	FOREIGN KEY (ticker)
	REFERENCES Holding(ticker)
	ON DELETE CASCADE;`)

}

func Test_AddNewTrade(t *testing.T) {
	testHolding := Holding{Ticker: "NKE", AvgCostPrice: 125.5, TotalShares: 45.5}
	if err := repository.AddNewTrade(context.TODO(), "testuserid", &testHolding); err != nil {
		t.Errorf("Unable to add new trade. failed : %v", err)
	}

}
