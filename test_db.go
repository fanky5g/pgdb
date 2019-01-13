package database

import (
	"strings"

	"github.com/flavioribeiro/gonfig"
)

func GetTestDbConfig() (gonfig.Gonfig, error) {
	reader := strings.NewReader(`
		{
			"dbHost": "localhost:5432",
			"dbName": "test",
			"dbUser": "test",
			"dbPass": "test"
		}
	`)

	return gonfig.FromJson(reader)
}
