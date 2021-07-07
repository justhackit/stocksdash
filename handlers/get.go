package handlers

import (
	"net/http"

	jsonutils "github.com/cloudlifter/go-utils/json"
)

func (sh *StocksdashHander) Stocksdashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey{}).(string)
	holdings, err := sh.repo.GetHoldings(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sh.logger.Error("Error while fetching holdings", "userid", userID, "error", err)
		jsonutils.ToJSON(&GenericResponse{
			Status:  false,
			Message: "hello," + userID + "! Oops, something went wrong.",
		}, w)
	} else {
		allHoldings := []DashboardAPIResponse{}
		holdedTickers := make([]string, len(*holdings))
		for _, holding := range *holdings {
			holdedTickers = append(holdedTickers, holding.Ticker)
		}
		currentPrices, _ := sh.repo.GetCurrentPrice(r.Context(), holdedTickers)
		for _, holding := range *holdings {
			thisHolding := DashboardAPIResponse{Ticker: holding.Ticker, AvgCostPrice: holding.AvgCostPrice, TotalShares: holding.TotalShares}
			thisHolding.CurrentPrice = currentPrices[thisHolding.Ticker]
			thisHolding.ProfitLoss = (thisHolding.CurrentPrice - thisHolding.AvgCostPrice) * thisHolding.TotalShares
			thisHolding.ProfitLossPerc = (thisHolding.CurrentPrice - thisHolding.AvgCostPrice) / thisHolding.AvgCostPrice * 100
			allHoldings = append(allHoldings, thisHolding)
		}
		w.WriteHeader(http.StatusOK)
		jsonutils.ToJSON(&allHoldings, w)
	}
}
