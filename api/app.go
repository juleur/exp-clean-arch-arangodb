package main

import (
	"context"
	"encoding/json"
	"exp-clean-arch-arangodb/api/routes"
	"exp-clean-arch-arangodb/pkg/activity"
	"exp-clean-arch-arangodb/pkg/city"
	"exp-clean-arch-arangodb/pkg/entities"
	"exp-clean-arch-arangodb/pkg/session"
	"exp-clean-arch-arangodb/pkg/user"
	"io/ioutil"
	"log"
	"os"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ArangoCollection struct {
	Type arangodb.CollectionType
	Name string
}

var COLLECTIONS = []ArangoCollection{
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "aliments",
	},
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "categories",
	},
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "communes",
	},
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "messengers",
	},
	{
		Type: arangodb.CollectionTypeEdge,
		Name: "relUsersAliments",
	},
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "sessions",
	},
	{
		Type: arangodb.CollectionTypeDocument,
		Name: "users",
	},
}

var CATEGORIES = []string{
	"fruit",
	"légume",
	"viande",
	"laitier",
	"oeuf",
}

func init() {
	arangoDB := arangoDBConnection()
	ctx := context.Background()

	for _, c := range COLLECTIONS {
		options := &arangodb.CreateCollectionOptions{
			Type: c.Type,
		}

		if found := arangoCollectionExists(arangoDB, c.Name); found {
			continue
		}

		_, err := arangoDB.CreateCollection(ctx, c.Name, options)
		if err != nil {
			log.Fatalln(err)
		}
	}
	// ajout des catégories
	if found := arangoCollectionExists(arangoDB, "categories"); found {
		return
	}

	col, err := arangoDB.Collection(ctx, "categories")
	if err != nil {
		log.Fatalln(err)
	}

	for _, cat := range CATEGORIES {
		doc := struct {
			Nom string `json:"nom"`
		}{
			Nom: cat,
		}
		_, err := col.CreateDocument(ctx, doc)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// ajout des communes
	if found := arangoCollectionExists(arangoDB, "communes"); found {
		return
	}

	col, err = arangoDB.Collection(ctx, "communes")
	if err != nil {
		log.Fatalln(err)
	}

	cities := openCommuneJSONFile()

	log.Printf("Ajout des %d communes Françaises en-cours...\n", len(*cities))
	for _, city := range *cities {
		_, err := col.CreateDocument(ctx, city)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// https://www.arangodb.com/docs/stable/analyzers.html
	// ne pas oublier de créer un analyzer trigram pour les communes
	// analyzers.save("ctrigram", "ngram", { min: 3, max: 3, preserveOriginal: true, streamType: "utf8" }, ["frequency", "norm", "position"]);
	// ne pas oublier la création des views pour les search
}

func main() {
	arangoDB := arangoDBConnection()

	activityRepo := activity.NewRepo(arangoDB)
	activiyService := activity.NewService(activityRepo)

	cityRepo := city.NewRepo(arangoDB)
	cityService := city.NewService(cityRepo)

	sessionRepo := session.NewRepo(arangoDB)
	sessionService := session.NewService(sessionRepo)

	userRepo := user.NewRepo(arangoDB)
	userService := user.NewService(userRepo)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Cookie, Authorization",
	}))

	api := app.Group("/api")

	routes.UserRouter(api, userService, sessionService, cityService, activiyService)
	routes.CityRouter(api, sessionService, cityService)

	_ = app.Listen(":2345")
}

func openCommuneJSONFile() *[]entities.CityCollection {
	jsonFile, err := os.Open("./../storages/communes/communes.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	cities := []entities.CityCollection{}
	err = json.Unmarshal(byteValue, &cities)
	if err != nil {
		log.Fatalln(err)
	}

	return &cities
}

func arangoCollectionExists(arangoDB arangodb.Database, collectionName string) bool {
	found, err := arangoDB.CollectionExists(context.Background(), collectionName)
	if err != nil {
		log.Fatalln(err)
	}
	return found
}

func arangoDBConnection() arangodb.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	if err != nil {
		log.Fatalln(err)
	}

	c, err := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatalln(err)
	}

	db, err := c.Database(context.Background(), "_system")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
