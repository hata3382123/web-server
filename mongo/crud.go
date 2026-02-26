package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			fmt.Println(failedEvent.Command)
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			fmt.Println(succeededEvent.Command)
		},
	}
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor)
	client, _ := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
}
