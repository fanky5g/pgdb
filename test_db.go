package database

import (
	"fmt"
	"strings"

	"os"

	"github.com/flavioribeiro/gonfig"
)

func getTestDbConfig() (gonfig.Gonfig, error) {
	dbUser := os.Getenv("DBUser")
	dbPass := os.Getenv("DBPass")

	reader := strings.NewReader(fmt.Sprintf(`
		{
			"dbHost": "/cloudsql/winst360:us-central1:winst-zero",
			"dbName": "postgres",
			"dbUser": "%s",
			"dbPass": "%s"
		}
	`, dbUser, dbPass))

	return gonfig.FromJson(reader)
}
