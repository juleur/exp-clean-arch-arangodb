package city

import (
	"context"
	"testing"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func TestFindCity(t *testing.T) {
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	cityRepo := NewRepo(arangoDB)
	cityService := NewService(cityRepo)

	cityName := "tremblay"
	citiesName, _ := cityService.FindCity(cityName)

	t.Log(citiesName)
}
