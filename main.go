package main

import (
	"context"
	"fmt"
	"os"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
	logutils "github.com/cloudlifter/go-utils/logs"
	configutils "github.com/justhackit/stocksdash/config"
)

func main() {
	logger := logutils.NewLogger()
	configFilepath := os.Args[1:][0]
	configs := configutils.NewConfigurationsFromFile(configFilepath, logger)
	logger.Info("config", configs)

	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "c2vuokiad3ifkigc36hg", // Replace this
	})
	//Stock candles
	stockCandles, _, err := finnhubClient.StockCandles(auth, "AAPL", "D", 1590988249, 1591852249, nil)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("%+v\n", stockCandles)

}
