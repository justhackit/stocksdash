package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	utils "github.com/cloudlifter/go-utils/timeutils"

	"github.com/hashicorp/go-hclog"
	"github.com/justhackit/stockedash/datastore"
)

type PriceHistoryResponse struct {
	Candles []struct {
		Open     float64 `json:"open"`
		High     float64 `json:"high"`
		Low      float64 `json:"low"`
		Close    float64 `json:"close"`
		Volume   int     `json:"volume"`
		Datetime int64   `json:"datetime"`
	} `json:"candles"`
	Symbol string `json:"symbol"`
	Empty  bool   `json:"empty"`
}

type TDAmeritradeAPI struct {
	apiToken   string
	logger     hclog.Logger
	repository datastore.Repository
}

func NewTDAmeritradeService(apitoken string, repo datastore.Repository, log hclog.Logger) *TDAmeritradeAPI {
	return &TDAmeritradeAPI{apiToken: apitoken, repository: repo, logger: log}
}

func (tdapi *TDAmeritradeAPI) GetHistoricalQuotes(ctx context.Context, ticker string, from time.Time, to time.Time) (*[]datastore.StockPrices, error) {
	defer utils.TimeTaken("GetHistoricalQuotes", tdapi.logger)()

	endpoint := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/%s/pricehistory?apikey=%s&periodType=year&period=1&frequencyType=daily&frequency=1&startDate=%d&endDate=%d&needExtendedHoursData=true",
		ticker, tdapi.apiToken, from.Unix()*1000, to.Unix()*1000)
	tdapi.logger.Info("TDAmeritradeEndpoint", "url", endpoint, "from", from.String(), "to", to.String())
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var quotes PriceHistoryResponse
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return nil, err
	}
	if quotes.Symbol != ticker {
		return nil, fmt.Errorf("i asked for ticker %s, but got %s", ticker, quotes.Symbol)
	}
	historicalPrices := make([]datastore.StockPrices, 0, len(quotes.Candles))
	for _, histPriceResp := range quotes.Candles {
		tmp := datastore.StockPrices{}
		tmp.Ticker = quotes.Symbol
		//Adding two hours because Ameritrade's Datetime has epoch as of 05 AM UTC for that trading day
		tmp.Date = time.Unix(histPriceResp.Datetime/1000, 0).Add(time.Hour * time.Duration(2))
		tmp.Open = histPriceResp.Open
		tmp.High = histPriceResp.High
		tmp.Low = histPriceResp.Low
		tmp.Close = histPriceResp.Close
		historicalPrices = append(historicalPrices, tmp)

	}
	return &historicalPrices, nil
}

func (tdapi *TDAmeritradeAPI) SaveHistoricalQuotes(ctx context.Context, ticker string, from time.Time, to time.Time) error {
	historicalPrices, err := tdapi.GetHistoricalQuotes(ctx, ticker, from, to)
	if err != nil {
		return err
	}
	tdapi.logger.Info("Saving Historical Price for DB", "ticker", ticker, "from", from, "to", to)
	if err := tdapi.repository.AddBatchHistorical(ctx, historicalPrices); err != nil {
		return err
	}
	return nil
}
