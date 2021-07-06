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
	Ticker                string  `json:"ticker"`
	AvgCostPrice          float64 `json:"avgCostPrice"`
	TotalShares           float64 `json:"totalShares"`
	ProfitLoss            float64 `json:"profitLoss"`
	ProfitLossPerc        float64 `json:"profitLossPerc"`
	TodaysChange          float64 `json:"todaysChange"`
	LastTwoDaysChange     float64 `json:"lastTwoDaysChange"`
	LastThreeDaysChange   float64 `json:"lastThreeDaysChange"`
	LastWeeksChange       float64 `json:"lastWeeksChange"`
	LastTwoWeeksChange    float64 `json:"lastTwoWeeksChange"`
	LastThreeWeeksChange  float64 `json:"lastThreeWeeksChange"`
	LastMonthsChange      float64 `json:"lastMonthsChange"`
	LastTwoMonthsChange   float64 `json:"lastTwoMonthsChange"`
	LastThreeMonthsChange float64 `json:"lastThreeMonthsChange"`
	LastSixMonthsChange   float64 `json:"lastSixMonthsChange"`
	LastYearsChange       float64 `json:"lastYearsChange"`
}
