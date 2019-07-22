package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPostgresConnection(t *testing.T) {
	config, err := getTestDbConfig()
	if assert.NoError(t, err) {
		db, err := GetPostgresConnectionWithDriver("postgres", config)
		if assert.NoError(t, err) {
			t.Log("pass")
			t.Log(db)
		}
	}
}
