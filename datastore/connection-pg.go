package datastore

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	utils "github.com/justhackit/stockedash/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConnection creates the connection to the database
func NewConnection(config *utils.Configurations, logger hclog.Logger) (*gorm.DB, error) {

	var conn string

	if config.Database.DBConn != "" {
		conn = config.Database.DBConn
	} else {
		host := config.Database.DBHost
		port := config.Database.DBPort
		user := config.Database.DBUser
		dbName := config.Database.DBName
		password := config.Database.DBPass
		schema := config.Database.DBSchema
		conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s search_path=%s sslmode=disable", host, port, user, dbName, password, schema)
	}
	logger.Debug("connection string", conn)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
