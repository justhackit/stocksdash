package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudlifter/go-utils/comms"
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
	tdameritrade = NewTDAmeritradeService(os.Getenv("TDAMERITRADE_TOKEN"), repo, logger)
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

func Test_GetCurrentQuote(t *testing.T) {
	testTickers := []string{"IBM", "NKE"}
	if currQuotes, err := tdameritrade.GetCurrentQuote(context.TODO(), testTickers); err != nil {
		t.Errorf("unable to get current quotes :%v", err)
	} else {
		for _, quote := range *currQuotes {
			fmt.Printf("%#v\n", quote)
		}
	}
}

func Test_SaveCurrentQuote(t *testing.T) {
	testTickers := []string{"IBM", "NKE"}
	comms.SendPushNotification("test", "test")
	for range time.Tick(time.Second * 5) {
		fmt.Printf("\n\n=========================\nRefreshing at %s\n", time.Now())
		if err := tdameritrade.SaveCurrentQuote(context.TODO(), testTickers); err != nil {
			t.Errorf("unable to save current quotes : %v", err)
		}
	}

}
