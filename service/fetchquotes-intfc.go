package service

import (
	"context"
	"time"

	"github.com/justhackit/stockedash/datastore"
)

type FetchQuotes interface {
	GetHistoricalQuotes(ctx context.Context, ticker string, from time.Time, to time.Time) (*[]datastore.StockPrices, error)
}
