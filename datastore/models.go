package datastore

import (
	"time"

	"gorm.io/gorm"
)

// User is the data type for user object
type Holding struct {
	gorm.Model
	UserId       string  `json:"userId" validate:"required" gorm:"column:userId;not null;primary_key"`
	Ticker       string  `json:"ticker" validate:"required" gorm:"column:ticker;not null;primary_key"`
	AvgCostPrice float64 `json:"avgCostPrice" validate:"required" gorm:"column:avgCostPrice;not null"`
	TotalShares  float64 `json:"totalShares" validate:"required" gorm:"column:avgCostPrice;not null"`
}

type HistoricalPrices struct {
	gorm.Model
	Ticker string    `gorm:"column:ticker;not null;primary_key"`
	Date   time.Time `gorm:"column:date;not null;primary_key"`
	Open   float64   `gorm:"column:open"`
	High   float64   `gorm:"column:high"`
	Low    float64   `gorm:"column:low"`
	Close  float64   `gorm:"column:close;not null"`
}
