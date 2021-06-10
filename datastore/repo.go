package datastore

import "context"

// Repository is an interface for the storage implementation of the auth service
type Repository interface {
	AddNewTrade(ctx context.Context, trade *Holding) error
	AddHistorical(ctx context.Context, tickers *HistoricalPrices) error
	//DeleteUser(ctx context.Context, email string, clientId string) error
}
