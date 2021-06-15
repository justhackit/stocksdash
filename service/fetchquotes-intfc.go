package service

import (
	"context"
	"time"

	"github.com/justhackit/stocksdash/datastore"
)

type FetchQuotes interface {
	GetHistoricalQuotes(ctx context.Context, ticker string, from time.Time, to time.Time) (*[]datastore.StockPrices, error)
	SaveHistoricalQuotes(ctx context.Context, ticker string, from time.Time, to time.Time) error
	GetCurrentQuote(ctx context.Context, tickers []string) error
	SaveCurrentQuote(ctx context.Context, tickers []string) error
}
