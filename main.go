package main

import (
	"context"
	"fmt"
	"os"
	"time"

	logutils "github.com/cloudlifter/go-utils/logs"
	configutils "github.com/justhackit/stocksdash/config"
	"github.com/justhackit/stocksdash/datastore"
	"github.com/justhackit/stocksdash/service"
)

func main() {
	logger := logutils.NewLogger()
	configFilepath := os.Args[1:][0]
	configs := configutils.NewConfigurationsFromFile(configFilepath, logger)
	logger.Info("config", configs)
	db, _ := datastore.NewConnection(configs, logger)
	repo := datastore.NewPostgresRepository(db, logger)
	tdameritrade := service.NewTDAmeritradeService("GJHIDO67W7GDJHPGUAOC9CHUKNEMXGOM", repo, logger)
	testTickers := []string{"IBM", "NKE"}
	for range time.Tick(time.Second * 5) {
		fmt.Printf("\n\n=========================\nRefreshing at %s\n", time.Now())
		if err := tdameritrade.SaveCurrentQuote(context.TODO(), testTickers); err != nil {
			fmt.Printf("unable to save current quotes : %v\n", err)
			panic(err)
		}
	}

}
