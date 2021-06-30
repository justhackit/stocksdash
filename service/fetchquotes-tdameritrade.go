package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudlifter/go-utils/comms"
	utils "github.com/cloudlifter/go-utils/timeutils"

	"github.com/hashicorp/go-hclog"
	"github.com/justhackit/stocksdash/datastore"
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
	if err := tdapi.repository.AddBatchQuotes(ctx, historicalPrices); err != nil {
		return err
	}
	return nil
}

func (tdapi *TDAmeritradeAPI) GetCurrentQuote(ctx context.Context, tickers []string) (*[]datastore.StockPrices, error) {
	defer utils.TimeTaken("GetCurrentQuote", tdapi.logger)()
	apiToken := os.Getenv("TDAMERITRADE_TOKEN")
	queryParams := fmt.Sprintf("apikey=%s&symbol=%s", apiToken, strings.Join(tickers[:], ","))
	endpoint := fmt.Sprintf("https://api.tdameritrade.com/v1/marketdata/quotes?%s", queryParams)
	fmt.Printf("Endpoint : %s\n", endpoint)
	req, _ := http.NewRequest("GET", endpoint, nil)

	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		tdapi.logger.Error("Error while accessing TD AMeritrade API", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]interface{}
	byteResp, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(byteResp, &result)
	stockPrices := make([]datastore.StockPrices, 0, len(tickers))
	for ticker := range result {
		tmp := datastore.StockPrices{}
		tmp.Ticker = ticker
		//Adding two hours because Ameritrade's Datetime has epoch as of 05 AM UTC for that trading day
		currentUTCInst := time.Now().In(time.UTC)
		year, month, day := currentUTCInst.Date()
		dayAt7UTC := time.Date(year, month, day, 0, 0, 0, 0, currentUTCInst.Location()).Add(time.Hour * time.Duration(7))
		tmp.Date = dayAt7UTC
		tmp.Open = result[ticker]["openPrice"].(float64)
		tmp.High = result[ticker]["highPrice"].(float64)
		tmp.Low = result[ticker]["lowPrice"].(float64)
		tmp.Close = result[ticker]["lastPrice"].(float64)
		stockPrices = append(stockPrices, tmp)
	}
	return &stockPrices, nil
}

func (tdapi *TDAmeritradeAPI) SaveCurrentQuote(ctx context.Context, tickers []string) error {
	defer utils.TimeTaken("SaveCurrentQuote", tdapi.logger)()
	currentPrices, err := tdapi.GetCurrentQuote(ctx, tickers)
	if err != nil {
		return err
	}
	tdapi.logger.Info("Saving Current Prices in DB", "tickers", tickers)
	if err := tdapi.repository.AddBatchQuotes(ctx, currentPrices); err != nil {
		return err
	}
	return nil
}

func (tdapi *TDAmeritradeAPI) KeepRefreshingQuotes(ctx context.Context) error {
	tdapi.logger.Info("Quotes refresher beginning...")
	for range time.Tick(time.Second * 30) {
		loc := time.FixedZone("UTC-7", -7*60*60)
		now := time.Now().In(loc)
		hr, _, _ := now.Clock()
		week := now.Weekday()
		tradeWeekDays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
		isATradingDay := false
		for _, aDay := range tradeWeekDays {
			if aDay == week {
				isATradingDay = true
			}
		}
		if hr > 6 && hr < 16 && isATradingDay {
			fmt.Printf("\n\n=========================\nRefreshing at %s\n", time.Now())
			holdings, err := tdapi.repository.GetHoldings(ctx)
			if err != nil {
				tdapi.logger.Error("unable to read holdings", "error", err)
			}
			var allTickers = make([]string, 0, len(*holdings))
			for _, holding := range *holdings {
				allTickers = append(allTickers, holding.Ticker)
			}
			if err := tdapi.SaveCurrentQuote(ctx, allTickers); err != nil {
				fmt.Printf("unable to save current quotes : %v\n", err)
				comms.SendPushNotification("Error getting current quotes", err.Error())
				return err
			}
		} else {
			tdapi.logger.Info("Not a trading hour..", "week", week, "hour", hr, "curr_time", time.Now())
		}
	}

	return nil

}
