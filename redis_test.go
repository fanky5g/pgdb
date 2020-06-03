package database

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestGetRedisConnection(t *testing.T) {
// 	cfg := RedisConnection{
// 		Master: ":6379",
// 		Slave:  ":6380",
// 	}

// 	client, err := GetRedisClient(cfg)
// 	assert.NoError(t, err)
// 	assert.Equal(t, MASTER_SLAVE, client.Mode)

// 	conn := client.GetConn("GET")
// 	// reply, err := conn.Do("SET", "pid:user-1", []byte("hello world"))
// 	// if assert.NoError(t, err) {
// 	// 	t.Log(reply)
// 	// }

// 	reply, err := conn.Do("GET", "pid:user-1")
// 	if assert.NoError(t, err) {
// 		t.Log(ByteArrayToString(reply.([]uint8)))
// 	}
// }

// // ByteArrayToString converts byte array into string
// func ByteArrayToString(bs []uint8) string {
// 	b := make([]byte, len(bs))
// 	for i, v := range bs {
// 		b[i] = byte(v)
// 	}
// 	return string(b)
// }
