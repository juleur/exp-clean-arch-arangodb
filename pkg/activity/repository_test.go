package activity

import (
	"context"
	"exp-clean-arch-arangodb/utils"
	"fmt"
	"testing"
	"time"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func TestActivity(t *testing.T) {
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	c, _ := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	arangoDB, _ := c.Database(context.Background(), "_system")

	activityRepo := NewRepo(arangoDB)
	activityService := NewService(activityRepo)

	usersKey := []string{
		"17669",
		"17670",
		"17671",
		"182219",
		"18443",
		"18444",
		"18445",
		"188568",
		"195509",
	}
	for _, k := range usersKey {
		encoded, _ := utils.EncodeString(k)
		token := utils.TokenGenerator(22)
		sessionID := fmt.Sprintf("%su%s", encoded, token)
		activityService.LastUserActivity(sessionID)
		t.Log(sessionID, encoded, k)
	}
	// simule infinite loop
	time.Sleep(5 * time.Second)
}
