package aliment

import (
	"context"
	"exp-clean-arch-arangodb/pkg/entities"
	"testing"
	"time"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	alimentRepo := NewRepo(arangoDB)
	alimentService := NewService(alimentRepo)

	categories, _ := alimentService.GetCategories()
	t.Log(categories)
}

func TestCreateAliment(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	alimentRepo := NewRepo(arangoDB)
	alimentService := NewService(alimentRepo)

	userID := "users/17670"
	alimentCollection := entities.AlimentCollection{
		Categorie:      "categories/20860",
		Nom:            "pomme",
		Variete:        "reinette clochard",
		SystemeEchange: []int{1, 2},
		Prix:           3.1,
		UniteMesure:    1,
		Stock:          2,
		AddedAt:        time.Now(),
	}
	err := alimentService.Create(userID, alimentCollection)
	assert.Nil(err)
}

func TestDeleteAliment(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	alimentRepo := NewRepo(arangoDB)
	alimentService := NewService(alimentRepo)

	userID := "users/17670"
	alimentID := "aliments/260847"
	err := alimentService.Delete(userID, alimentID)
	assert.Nil(err)
}

func TestFetchCategorie(t *testing.T) {
	assert := assert.New(t)
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	alimentRepo := NewRepo(arangoDB)

	categorieName := "fruit"
	catResult, err := alimentRepo.fetchCategorie(categorieName)

	assert.Nil(err)
	assert.Equal(categorieName, catResult)
}
