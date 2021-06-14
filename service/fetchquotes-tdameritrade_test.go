package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	logutils "github.com/cloudlifter/go-utils/logs"
	"github.com/justhackit/stocksdash/config"
	"github.com/justhackit/stocksdash/datastore"
)

var tdameritrade *TDAmeritradeAPI

func init() {
	logger := logutils.NewLoggerWithName("stocksdash", "DEBUG")
	configs := config.NewConfigurationsFromFile("../test-config.yaml", logger)
	db, _ := datastore.NewConnection(configs, logger)
	repo := datastore.NewPostgresRepository(db, logger)
	tdameritrade = NewTDAmeritradeService("GJHIDO67W7GDJHPGUAOC9CHUKNEMXGOM", repo, logger)
}

func Test_GetHistoricalQuotes(t *testing.T) {
	testTicker := "NKE"
	from := time.Date(2021, time.June, 10, 5, 1, 0, 0, time.UTC)
	to := time.Date(2021, time.June, 11, 5, 1, 0, 0, time.UTC)
	historicals, err := tdameritrade.GetHistoricalQuotes(context.TODO(), testTicker, from, to)
	if err != nil {
		t.Errorf("Unable to fetch historical prices: %v", err)
	} else {
		for _, aquote := range *historicals {
			fmt.Printf("Ticker :%s, Day : %s Close : %f\n", aquote.Ticker, aquote.Date.String(), aquote.Close)
		}
	}

}

func Test_SaveHistoricalQuotes(t *testing.T) {
	testTicker := "IBM"
	from := time.Date(2010, time.June, 10, 5, 1, 0, 0, time.UTC)
	to := time.Date(2021, time.June, 11, 5, 1, 0, 0, time.UTC)
	if err := tdameritrade.SaveHistoricalQuotes(context.TODO(), testTicker, from, to); err != nil {
		t.Errorf("unable to save historical quotes : %v", err)
	}
}
