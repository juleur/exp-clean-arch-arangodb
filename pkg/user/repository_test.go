package user

import (
	"context"
	"exp-clean-arch-arangodb/pkg/entities"
	"testing"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	userRepo := NewRepo(arangoDB)
	userService := NewService(userRepo)

	newUser := entities.NewUserBody{
		Username: "cel45",
		Email:    "cel-test@gmail.com",
		Password: "12345",
	}

	username, err := userService.Create(newUser)
	assert.Nil(err)
	assert.Equal(newUser.Username, username)
}
