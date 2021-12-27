package datastore

import (
	"time"
)

// User is the data type for user object
type Holding struct {
	UserId       string  `json:"userId" validate:"required" gorm:"column:user_id;not null;primary_key"`
	Ticker       string  `json:"ticker" validate:"required" gorm:"column:ticker;not null;primary_key"`
	AvgCostPrice float64 `json:"avgCostPrice" validate:"required" gorm:"column:avg_cost_price;not null"`
	TotalShares  float64 `json:"totalShares" validate:"required" gorm:"column:total_shares;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

type StockPrices struct {
	Ticker    string    `gorm:"column:ticker;not null;primary_key"`
	Date      time.Time `gorm:"column:date;not null;primary_key"`
	Open      float64   `gorm:"column:open"`
	High      float64   `gorm:"column:high"`
	Low       float64   `gorm:"column:low"`
	Close     float64   `gorm:"column:close;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
