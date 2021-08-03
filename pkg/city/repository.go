package city

import (
	"context"
	"errors"
	"exp-clean-arch-arangodb/pkg/entities"

	arangodb "github.com/arangodb/go-driver"
)

const DefaultError = "Une erreur est survenue"

type Repository interface {
	GetCity(cityName string) ([]string, error)
	CheckIfCityNameExists(cityName string) (*entities.CityCollection, error)
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangodb arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangodb,
	}
}

func (r *repository) GetCity(cityName string) ([]string, error) {
	query := `
		FOR c IN communesView
			SEARCH NGRAM_MATCH(
				c.nom,
				@city,
				0.45,
				"ctrigram"
			)
			SORT BM25(c) DESC
			LIMIT 15
			RETURN c.nom
	`
	bindVars := map[string]interface{}{
		"city": cityName,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return []string{}, errors.New(DefaultError)
	}
	defer cursor.Close()

	citiesName := []string{}
	for {
		var doc string
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return []string{}, errors.New(DefaultError)
		}
		citiesName = append(citiesName, doc)
	}
	return citiesName, nil
}

func (r *repository) CheckIfCityNameExists(cityName string) (*entities.CityCollection, error) {
	query := `
		FOR c IN communes
			FILTER c.nom == @city
			RETURN c
	`
	bindVars := map[string]interface{}{
		"city": cityName,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New(DefaultError)
	}
	defer cursor.Close()

	city := entities.CityCollection{}
	if _, err := cursor.ReadDocument(context.Background(), &city); err != nil {
		return nil, errors.New(DefaultError)
	}

	return &city, nil
}
