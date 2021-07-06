package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudlifter/go-utils/comms"
	logutils "github.com/cloudlifter/go-utils/logs"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	configutils "github.com/justhackit/stocksdash/config"
	"github.com/justhackit/stocksdash/datastore"
	"github.com/justhackit/stocksdash/handlers"
	"github.com/justhackit/stocksdash/service"
)

func main() {
	logger := logutils.NewLoggerWithName("stocksdash", "INFO")
	configFilepath := os.Args[1:][0]
	logger.Info("logfile path", configFilepath)
	configs := configutils.NewConfigurationsFromFile(configFilepath, logger)
	logger.Info("config", configs)
	db, err := datastore.NewConnection(configs, logger)
	if err != nil {
		logger.Error("unable to initialize DB")
		panic(err)
	}
	//Setup db if this is running for the first time
	db.Debug().AutoMigrate(&datastore.Holding{}, &datastore.StockPrices{})
	repo := datastore.NewPostgresRepository(db, logger)
	go func() {
		tdAmeritradeToken := os.Getenv("TDAMERITRADE_TOKEN")
		tdameritrade := service.NewTDAmeritradeService(tdAmeritradeToken, repo, logger)
		maxRetriesAllowed := 30
		for timeoutsSoFar := 0; timeoutsSoFar < maxRetriesAllowed; timeoutsSoFar++ {
			tdameritrade.KeepRefreshingQuotes(context.TODO())
			logger.Error("Error while refreshing quotes. Retrying after after waiting for 5 mins... ", "timeOutsSoFar", timeoutsSoFar)
			time.Sleep(5 * time.Minute)
		}
		comms.SendPushNotification("main.go", fmt.Sprintf("Max retries exhausted : %d", maxRetriesAllowed))
	}()

	apiService := service.NewStockdashSvc(logger, configs)
	serviceHandler := handlers.NewStocksdashHandler(logger, configs, repo, apiService)
	sm := mux.NewRouter().PathPrefix("/stocksdash").Subrouter()
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/dashboard", serviceHandler.Stocksdashboard)
	getR.Use(serviceHandler.MiddlewareValidateAccessToken)

	// create a server
	svr := http.Server{
		Addr:         configs.ServerAddress,
		Handler:      sm,
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		logger.Info("starting the server at port", configs.ServerAddress)

		err := svr.ListenAndServe()
		if err != nil {
			logger.Error("could not start the server", "error", err)
			os.Exit(1)
		}
	}()

	// look for interrupts for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	logger.Info("shutting down the server", "received signal", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	svr.Shutdown(ctx)

}
