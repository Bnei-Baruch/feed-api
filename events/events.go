package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	eventsMutex sync.Mutex
	events      []Event
	DebugMode   bool
)

func ReadAndClearEvents() []Event {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	ret := make([]Event, len(events))
	copy(ret, events)
	events = []Event(nil)

	return ret
}

func AddEvent(e Event) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	events = append(events, e)
}

// Event data struct for unmarshaling data from nats
type Event struct {
	ID                  string                 `json:"id"`
	Type                string                 `json:"type"`
	ReplicationLocation string                 `json:"rloc"`
	Payload             map[string]interface{} `json:"payload"`
}

// Nats
type EventListener struct {
	nc         *nats.Conn
	js         jetstream.JetStream
	consumer   jetstream.Consumer
	consumeCtx jetstream.ConsumeContext
}

// Returns shut down function.
func RunListener() func() {
	el := new(EventListener)

	DebugMode = viper.GetString("nats.mode") == "debug"
	natsURL := viper.GetString("nats.url")
	log.Infof("Initialize connection to nats debug=%t %s", DebugMode, natsURL)

	var err error
	el.nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Errorf("nats.Connect: %w", err)
		utils.Must(err)
	}

	el.js, err = jetstream.New(el.nc)
	if err != nil {
		log.Errorf("jetstream.New: %w", err)
		utils.Must(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	el.consumer, err = el.js.CreateOrUpdateConsumer(ctx, "MDB", jetstream.ConsumerConfig{
		Name:        "Feed-API",
		Durable:     "Feed-API",
		Description: "Events listener of MDB",
	})
	if err != nil {
		log.Errorf("jetstream.CreateOrUpdateConsumer: %w", err)
		utils.Must(err)
	}

	el.consumeCtx, err = el.consumer.Consume(el.handleMessage)
	if err != nil {
		log.Errorf("jetstream consumer.Consume: %w", err)
		utils.Must(err)
	}

	return func() { el.Close() }
}

func (el *EventListener) Close() {
	el.consumeCtx.Stop()
	el.nc.Close()
}

func (el *EventListener) handleMessage(msg jetstream.Msg) {
	log.Debugf("EventListener.handleMessage: %+v", msg.Data())

	var event Event
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		log.Errorf("EventListener.handleMessage json.Unmarshal: %w", err)
	}

	log.Debugf("EventListener.handleMessage event: %+v", event)
	AddEvent(event)

	if !DebugMode {
		msg.Ack()
	}
}
