package datastore

import (
	"context"
	"testing"
	"time"

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
	db.Debug().AutoMigrate(&Holding{}, &StockPrices{})

}

func Test_AddNewTrade(t *testing.T) {
	testHolding := Holding{Ticker: "NKE", AvgCostPrice: 125.5, TotalShares: 45.5}
	if err := repository.AddNewTrade(context.TODO(), "testuserid", &testHolding); err != nil {
		t.Errorf("Unable to add new trade. failed : %v", err)
	}
	//	all, _ := repository.GetHoldingsByUser(context.TODO(), "testuserid")
	//fmt.Printf("%#v", all)

}

func Test_GetHoldingsByUser(t *testing.T) {
	testUserId := "testuserid"
	if holdings, err := repository.GetHoldingsByUser(context.TODO(), testUserId); err != nil {
		t.Errorf("Error while getting holdings of an user : %v", err)
		t.Logf("All holdings: %#v", holdings)
		if len(*holdings) == 0 {
			t.Errorf("Got 0 records for the user %s", testUserId)
		}
	}

}

func Test_AddHistorical(t *testing.T) {
	testHistorical := &StockPrices{Ticker: "NKE", Date: time.Date(2020, time.April, 12, 20, 0, 0, 0, time.UTC), Open: 136.44, High: 147.95, Low: 134.67, Close: 141.47}
	if err := repository.AddHistorical(context.TODO(), testHistorical); err != nil {
		t.Errorf("Unable to add historical price: %v", err)
	}
}
