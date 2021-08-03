package aliment

import (
	"context"
	"errors"
	"exp-clean-arch-arangodb/pkg/entities"

	arangodb "github.com/arangodb/go-driver"
)

const DefaultError = "Une erreur est survenue"
const AlimentRequirements = "Vous devez enregistrer une commune et au moins, un moyen de communication"

type Repository interface {
	FetchCategories() ([]string, error)
	InsertAliment(userID string, aliment entities.AlimentCollection) error
	DeleteAliment(userID string, alimentID string) error
	fetchCategorie(categorieName string) (string, error)
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangodb arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangodb,
	}
}

func (r *repository) FetchCategories() ([]string, error) {
	query := "FOR cat IN categories RETURN cat.nom"
	cursor, err := r.ArangoDB.Query(context.Background(), query, nil)
	if err != nil {
		return []string{}, errors.New(DefaultError)
	}
	defer cursor.Close()

	categoriesName := []string{}
	for {
		var doc string
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return []string{}, errors.New(DefaultError)
		}
		categoriesName = append(categoriesName, doc)
	}
	return categoriesName, nil
}

func (r *repository) InsertAliment(userID string, aliment entities.AlimentCollection) error {
	// vérifie si l'utilisateur a bien enregistré une commune
	// et au moins, une messagerie pour le contacter
	hasAlimentRequirements := r.hasUserCityAndMessengers(userID)
	if !hasAlimentRequirements {
		return errors.New(AlimentRequirements)
	}

	col, err := r.ArangoDB.Collection(context.Background(), "aliments")
	if err != nil {
		return errors.New(DefaultError)
	}

	meta, err := col.CreateDocument(context.Background(), aliment)
	if err != nil {
		return errors.New(DefaultError)
	}

	relUsersAliments := entities.RelUsersAliments{
		From: userID,
		To:   meta.ID.String(),
	}

	err = r.insertRelationshipUsersAliments(relUsersAliments)
	if err != nil {
		// supprime l'aliment si la relation a échouée
		_, err := col.RemoveDocument(context.Background(), meta.ID.String())
		if err != nil {
			return errors.New(DefaultError)
		}

		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) DeleteAliment(userID string, alimentID string) error {
	col, err := r.ArangoDB.Collection(context.Background(), "aliments")
	if err != nil {
		return errors.New(DefaultError)
	}

	_, err = col.RemoveDocument(context.Background(), alimentID)
	if err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	query := `
		FOR rua IN relUsersAliments
			FILTER rua._from == @userID
			FILTER rua._to == @alimentID
			REMOVE { _key: rua._key } IN relUsersAliments
	`
	bindVars := map[string]interface{}{
		"userID":    userID,
		"alimentID": alimentID,
	}

	if _, err = r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) insertRelationshipUsersAliments(relUsersAliments entities.RelUsersAliments) error {
	col, err := r.ArangoDB.Collection(context.Background(), "relUsersAliments")
	if err != nil {
		return errors.New(DefaultError)
	}

	_, err = col.CreateDocument(context.Background(), relUsersAliments)
	if err != nil {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) fetchCategorie(categorieName string) (string, error) {
	query := "FOR cat IN categories FILTER cat.nom == @categorie RETURN cat.nom"
	bindVars := map[string]interface{}{
		"categorie": categorieName,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return "", errors.New(DefaultError)
	}
	defer cursor.Close()

	categoriesName := []string{}
	for {
		var doc string
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return "", errors.New(DefaultError)
		}
		categoriesName = append(categoriesName, doc)
	}

	return categoriesName[0], nil
}

func (r *repository) hasUserCityAndMessengers(userID string) bool {
	query := `
		FOR user IN users
			FILTER user.commune != NULL
			FILTER user.messenger != NULL
			FILTER user._id == @userID
			return user
	`
	bindVars := map[string]interface{}{
		"userID": userID,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return false
	}
	defer cursor.Close()

	meta, err := cursor.ReadDocument(context.Background(), &entities.ArangoResult{})
	if err != nil {
		return false
	}

	return !meta.ID.IsEmpty()
}
