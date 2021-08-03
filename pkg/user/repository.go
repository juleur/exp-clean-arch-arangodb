package user

import (
	"context"
	"errors"
	"exp-clean-arch-arangodb/pkg/entities"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	arangodb "github.com/arangodb/go-driver"
)

const EmailTaken = "Cet email est déjà utilisé"
const UsernameTaken = "Ce nom d'utilisateur est déjà utilisé"
const BadCredentials = "Ces informations ne correspondent à aucun utilisateur"
const DeleteUserErr = "Une erreur s'est produite lors de la suppression de votre compte"
const DefaultError = "Une erreur est survenue"

type Repository interface {
	InsertUser(newUser entities.NewUserBody) (string, error)
	FetchUserByEmail(email string) (*entities.UserDoc, error)
	UpdateLastSeenUser(userID string) error
	InsertCityUser(userID string, cityName string) error
	DeleteCityUser(userID string) error
	DeleteTotallyUser(userID string) error
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangodb arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangodb,
	}
}

func (r *repository) InsertUser(newUser entities.NewUserBody) (string, error) {
	query := "FOR user IN users FILTER user.email == @email RETURN user"
	bindVars := map[string]interface{}{
		"email": newUser.Email,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return "", errors.New(DefaultError)
	}

	meta, err := cursor.ReadDocument(context.Background(), &entities.ArangoResult{})
	if err != nil && !arangodb.IsNoMoreDocuments(err) {
		return "", errors.New(DefaultError)
	}

	if meta.Key != "" {
		return "", errors.New(EmailTaken)
	}

	query = "FOR user IN users FILTER user.pseudo == @pseudo RETURN user"
	bindVars = map[string]interface{}{
		"pseudo": newUser.Username,
	}
	cursor, err = r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return "", errors.New(DefaultError)
	}
	defer cursor.Close()

	meta, err = cursor.ReadDocument(context.Background(), &entities.ArangoResult{})
	if err != nil && !arangodb.IsNoMoreDocuments(err) {
		return "", errors.New(DefaultError)
	}

	if meta.Key != "" {
		return "", errors.New(UsernameTaken)
	}

	hashedPWD, err := argon2id.CreateHash(newUser.Password, argon2id.DefaultParams)
	if err != nil {
		return "", errors.New(DefaultError)
	}

	newUserDoc := entities.UserDoc{
		Pseudo:       newUser.Username,
		Email:        newUser.Email,
		Hpwd:         hashedPWD,
		RegisteredAt: time.Now(),
	}
	col, err := r.ArangoDB.Collection(context.Background(), "users")
	if err != nil {
		return "", errors.New(DefaultError)
	}

	if _, err = col.CreateDocument(context.Background(), newUserDoc); err != nil {
		return "", errors.New(DefaultError)
	}

	return newUser.Username, nil
}

func (r *repository) FetchUserByEmail(email string) (*entities.UserDoc, error) {
	query := "FOR user IN users FILTER user.email == @email return user"
	bindVars := map[string]interface{}{
		"email": email,
	}

	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil && !arangodb.IsNotFound(err) {
		return nil, errors.New(DefaultError)
	}
	defer cursor.Close()

	userDoc := entities.UserDoc{}
	if _, err := cursor.ReadDocument(context.Background(), &userDoc); err != nil {
		return nil, errors.New(DefaultError)
	}

	return &userDoc, nil
}

func (r *repository) UpdateLastSeenUser(userID string) error {
	query := `UPDATE "@userID" WITH { lastSeenAt: @timestamp } IN users`
	bindVars := map[string]interface{}{
		"userID":    userID,
		"timestamp": time.Now(),
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) InsertCityUser(userID string, cityID string) error {
	query := `UPDATE "@userID" WITH { commune: @cityID } IN users`
	bindVars := map[string]interface{}{
		"userID": userID,
		"cityID": cityID,
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) DeleteCityUser(userID string) error {
	query := `UPDATE "@userID" WITH { commune: NULL } IN users`
	bindVars := map[string]interface{}{
		"userID": userID,
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) DeleteTotallyUser(userID string) error {
	col, err := r.ArangoDB.Collection(context.Background(), "users")
	if err != nil {
		return err
	}

	if _, err = col.RemoveDocument(context.Background(), userID); err != nil {
		return err
	}

	if err = r.deleteRelationUsersAliments(userID); err != nil {
		return errors.New(DefaultError)
	}

	if err = r.deleteUserMessenger(userID); err != nil {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) deleteRelationUsersAliments(userID string) error {
	query := `
		FOR rua IN relUsersAliments
			FILTER rua._from == @userID
			REMOVE { _key: rua._key } IN relUsersAliments
	`
	bindVars := map[string]interface{}{
		"userID": fmt.Sprintf("users/%s", userID),
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *repository) deleteUserMessenger(userID string) error {
	query := `
		FOR user IN users
			FILTER user.messenger != NULL
			FILTER user._key == @userID
			FOR messenger IN messengers
				FILTER messenger._id == user.messenger
				REMOVE { _key: messenger._key } IN messengers
	`
	bindVars := map[string]interface{}{
		"userID": userID,
	}

	if _, err := r.ArangoDB.Query(context.Background(), query, bindVars); err != nil && !arangodb.IsNotFound(err) {
		return errors.New(DefaultError)
	}

	return nil
}
