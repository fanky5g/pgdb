package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRedisConnection(t *testing.T) {
	cfg := RedisConnection{
		Master: ":6379",
	}

	client, err := GetRedisClient(cfg)
	if assert.NoError(t, err) {
		if assert.Equal(t, MASTER_SLAVE, client.Mode) {
			t.Log("Successfully created redis client")
		}
	}
}
