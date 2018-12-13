package database

import (
	"strings"

	"github.com/flavioribeiro/gonfig"
)

func GetTestDbConfig() (gonfig.Gonfig, error) {
	reader := strings.NewReader(`
		{
			"dbHost": "",
			"dbName": "",
			"dbUser": "",
			"dbPass": ""
		}
	`)

	return gonfig.FromJson(reader)
}
