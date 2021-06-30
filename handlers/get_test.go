package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	logutils "github.com/cloudlifter/go-utils/logs"
	configutils "github.com/justhackit/stocksdash/config"
	"github.com/justhackit/stocksdash/datastore"
	"github.com/justhackit/stocksdash/service"
)

var stockHandler *StocksdashHander

func init() {
	logger := logutils.NewLoggerWithName("stocksdash", "INFO")
	configs := configutils.NewConfigurationsFromFile("../test-config.yaml", logger)
	db, _ := datastore.NewConnection(configs, logger)
	repo := datastore.NewPostgresRepository(db, logger)
	apiService := service.NewStockdashSvc(logger, configs)
	stockHandler = NewStocksdashHandler(logger, configs, repo, apiService)
}

func Test_Stocksdashboard(t *testing.T) {
	req, _ := http.NewRequest("GET", "/stocksdash/dashboard", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIzZTdjNTgxMS1lZGY1LTRhZDgtOTA3Ny0yZWI2MTdmYTIxMzQiLCJLZXlUeXBlIjoiYWNjZXNzIiwiZXhwIjoxNjI0OTI4NzQ5LCJpc3MiOiJjbG91ZGxpZnRlci5hdXRoLnNlcnZpY2UifQ.lIxV0n4hcs45won6JJ5BirkA1gUgSTv9lMF_u2fFY_ZedsPCDepQ3yYQJPxeI9ZzkoM_5S3KLUuZ9T98PSql-Q3nwNeOywlWHW0ZCyd45vWs4DQTYSJWAQzOISZycxpW_8ayStVwinbraZeO1JXqagmyeXsDovh5JxU5JuO9MUsgjeA5gWYrRBK50vPWuoFbbYYsYnxLjZBzzuZaA2EISEua7cNCc1erU-Q0s_he_lk071H3jCkDUkNfePI9h2Wxi2sAQgNeM_eNlePuspyINxN5HpAaEHfkzYm9Abe2CMgv9l6IEPGwwNkp-d6JjMdT8KgxCYmJNvStbW0g24mbNw")
	handler := http.HandlerFunc(stockHandler.Stocksdashboard)
	respRecord := httptest.NewRecorder()
	handler.ServeHTTP(respRecord, req)
	if status := respRecord.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
