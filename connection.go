package database

import (
	"fmt"

	"github.com/flavioribeiro/gonfig"
	"github.com/jinzhu/gorm"
	// brings postgresql into scope
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

// GetPostgresConnection returns postgres database connection
func GetPostgresConnection(config gonfig.Gonfig) (*gorm.DB, error) {
	dbName, err := config.GetString("dbName", "postgres")
	if err != nil {
		return nil, err
	}

	dbHost, err := config.GetString("dbHost", "")
	if err != nil {
		return nil, err
	}

	dbUser, err := config.GetString("dbUser", "")
	if err != nil {
		return nil, err
	}

	dbPass, err := config.GetString("dbPass", "")

	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPass)
	db, err := gorm.Open("cloudsqlpostgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
