package datastore

import (
	"context"
)

// Repository is an interface for the storage implementation of the auth service
type Repository interface {
	AddNewTrade(ctx context.Context, userId string, trade *Holding) error
	//DeleteATicker(ctx context.Context, ticker string, userId string) error //yet to implement
	GetHoldingsByUser(ctx context.Context, userId string) (*[]Holding, error)
	AddHistorical(ctx context.Context, tickers *StockPrices) error
	AddBatchHistorical(ctx context.Context, tickers *[]StockPrices) error
	//DeleteHistorical(ctx context.Context, ticker string, timeBefore ...time.Time) error        //yet to implement
	//GetEODPrice(ctx context.Context, ticker string, date time.Time) (*HistoricalPrices, error) //yet to implement
	//DeleteUser(ctx context.Context, email string, clientId string) error
}
