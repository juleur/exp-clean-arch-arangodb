package session

import (
	"context"
	"testing"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateSessionID(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	sessionRepo := NewRepo(arangoDB)
	sessionService := NewService(sessionRepo)

	userID := "17671"

	sessionID, err := sessionService.Set(userID)
	assert.Nil(err)
	t.Log(sessionID)
}

func TestGetSessionID(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	sessionRepo := NewRepo(arangoDB)
	sessionService := NewService(sessionRepo)

	sessionIDExists := "46gF2yVudmM6LLxvecVLDkTEbxDCO6"
	err := sessionService.Get(sessionIDExists)
	assert.Nil(err)

	sessionIDNotExist := "46gF2yVuomM9PLxvgcVLKkTEbxDCO"
	err = sessionService.Get(sessionIDNotExist)
	assert.Nil(err)
}

func TestRemoveSessionID(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	sessionRepo := NewRepo(arangoDB)
	sessionService := NewService(sessionRepo)

	sessionIDExists := "46gF2yVudmM6LLxvecVLDkTEbxDCO6"

	err := sessionService.Delete(sessionIDExists)
	assert.Nil(err)
}
