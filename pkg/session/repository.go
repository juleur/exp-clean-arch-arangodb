package session

import (
	"context"
	"errors"
	"exp-clean-arch-arangodb/pkg/entities"
	"exp-clean-arch-arangodb/utils"
	"fmt"
	"log"

	arangodb "github.com/arangodb/go-driver"
)

const SessionErr = "Une erreur de session s'est produite, veuillez vous authentifier une nouvelle fois"
const SessionNotFound = "Cette session est introuvable"

type Repository interface {
	GetUserSessionID(sessionID string) error
	PutUserInSession(userID string) (string, error)
	RemoveUserFromSession(sessionID string) error
}

type repository struct {
	SessionDB arangodb.Collection
}

func NewRepo(arangoDB arangodb.Database) Repository {
	sessionCollection, err := arangoDB.Collection(context.Background(), "sessions")
	if err != nil {
		log.Fatalln(err)
	}
	return &repository{
		SessionDB: sessionCollection,
	}
}

func (r *repository) GetUserSessionID(sessionID string) error {
	meta, err := r.SessionDB.ReadDocument(context.Background(), sessionID, &entities.ArangoResult{})
	if err != nil {
		return err
	}
	if meta.Key == "" {
		return errors.New(SessionNotFound)
	}
	return nil
}

func (r *repository) PutUserInSession(userID string) (string, error) {
	sessionID := utils.TokenGenerator(22)

	userIDEnc, err := utils.EncodeString(userID)
	if err != nil {
		return "", err
	}

	fullSessionID := fmt.Sprintf("%su%s", userIDEnc, sessionID)

	newSessionDoc := entities.SessionCollection{
		Key: fullSessionID,
	}
	if _, err = r.SessionDB.CreateDocument(context.Background(), newSessionDoc); err != nil {
		return "", err
	}

	return fullSessionID, nil
}

func (r *repository) RemoveUserFromSession(sessionID string) error {
	if _, err := r.SessionDB.RemoveDocument(context.Background(), sessionID); err != nil {
		if arangodb.IsNotFound(err) {
			return errors.New(SessionNotFound)
		}
		return errors.New(SessionErr)
	}
	return nil
}
