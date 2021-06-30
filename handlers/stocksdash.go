package handlers

import (
	"github.com/hashicorp/go-hclog"
	configutils "github.com/justhackit/stocksdash/config"
	"github.com/justhackit/stocksdash/datastore"
	"github.com/justhackit/stocksdash/service"
)

// UserIDKey is used as a key for storing the UserID in context at middleware
type UserIDKey struct{}

type StocksdashHander struct {
	logger        hclog.Logger
	configs       *configutils.Configurations
	repo          datastore.Repository
	stocksService service.Stocksdash
}

func NewStocksdashHandler(l hclog.Logger, c *configutils.Configurations, r datastore.Repository, svc service.Stocksdash) *StocksdashHander {
	return &StocksdashHander{logger: l, configs: c, repo: r, stocksService: svc}
}

// GenericResponse is the format of our response
type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DashboardAPIResponse struct {
	Ticker       string  `json:"ticker"`
	AvgCostPrice float64 `json:"avgCostPrice"`
	TotalShares  float64 `json:"totalShares"`
}
