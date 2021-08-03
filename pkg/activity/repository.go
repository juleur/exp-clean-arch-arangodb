package activity

import (
	"context"
	"exp-clean-arch-arangodb/utils"
	"strings"
	"time"

	arangodb "github.com/arangodb/go-driver"
)

type Repository interface {
	UpdateLastSeenUser(sessionID string)
}

type repository struct {
	activityQUpdate chan string
}

func NewRepo(arangodb arangodb.Database) Repository {
	activityQUpdate := make(chan string, 10)

	go func() {
		for sessionID := range activityQUpdate {
			split := strings.Split(sessionID, "u")

			key, err := utils.DecodeString(split[0])
			if err != nil {
				continue
			}

			query := `UPDATE @key WITH { lastSeenAt: @datetime } IN users`
			bindVars := map[string]interface{}{
				"key":      key,
				"datetime": time.Now(),
			}
			_, err = arangodb.Query(context.Background(), query, bindVars)
			if err != nil {
				continue
			}
		}
	}()

	return &repository{
		activityQUpdate: activityQUpdate,
	}
}

func (r *repository) UpdateLastSeenUser(sessionID string) {
	r.activityQUpdate <- sessionID
}
