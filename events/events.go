package events

import (
	"encoding/json"
	"runtime/debug"
	"sync"
	"time"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	eventsMutex  sync.Mutex
	events       []Data
	resyncNeeded bool
	DebugMode    bool
)

// ReadAndClearEvents also reports whether we reconnected to nats since the last call,
// events from the time we were disconnected may be lost.
func ReadAndClearEvents() ([]Data, bool) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	defer func() { events, resyncNeeded = nil, false }()
	return events, resyncNeeded
}

func AddEvent(d Data) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	events = append(events, d)
}

// Returns shut down function.
func RunListener() func() {
	log.Info("Initialize connection to nats")
	natsURL := viper.GetString("nats.url")
	natsClientID := viper.GetString("nats.client-id")
	natsClusterID := viper.GetString("nats.cluster-id")
	natsSubject := viper.GetString("nats.subject")
	DebugMode = viper.GetString("nats.mode") == "debug"

	var startOpt stan.SubscriptionOption
	if viper.GetBool("nats.durable") == true {
		startOpt = stan.DurableName(viper.GetString("nats.durable-name"))
	} else {
		startOpt = stan.DeliverAllAvailable()
	}

	var scMutex sync.Mutex
	var sc stan.Conn
	closed := false
	var connect func() error
	connect = func() error {
		conn, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL),
			stan.Pings(10, 6), // Connection is lost after ~1 minute without nats.
			stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
				log.Errorf("Nats connection lost: %v. Reconnecting.", reason)
				for err := connect(); err != nil; err = connect() {
					log.Errorf("Nats reconnect failed: %v", err)
					time.Sleep(10 * time.Second)
				}
				log.Info("Nats reconnected.")
				eventsMutex.Lock()
				resyncNeeded = true
				eventsMutex.Unlock()
			}))
		if err != nil {
			return err
		}
		log.Info("Subscribing to nats subject")
		if _, err = conn.Subscribe(natsSubject, msgHandler, startOpt, stan.SetManualAckMode()); err != nil {
			conn.Close()
			return err
		}
		scMutex.Lock()
		defer scMutex.Unlock()
		if closed {
			conn.Close()
			return nil
		}
		sc = conn
		return nil
	}
	utils.Must(connect())

	return func() {
		scMutex.Lock()
		defer scMutex.Unlock()
		closed = true
		sc.Close()
	}
}

// Data struct for unmarshaling data from nats
type Data struct {
	ID                  string                 `json:"id"`
	Type                string                 `json:"type"`
	ReplicationLocation string                 `json:"rloc"`
	Payload             map[string]interface{} `json:"payload"`
}

// msgHandler checks message type and calls "eventHandler"
func msgHandler(msg *stan.Msg) {
	// don't panic !
	defer func() {
		if rval := recover(); rval != nil {
			log.Errorf("msgHandler panic: %v while handling %v", rval, msg)
			debug.PrintStack()
		}
	}()

	var d Data
	err := json.Unmarshal(msg.Data, &d)
	if err != nil {
		log.Errorf("json.Unmarshal error: %s\n", err)
	}

	// Acknowledge the message
	if !DebugMode {
		msg.Ack()
	}

	log.Infof("Events Adding %+v", d)
	AddEvent(d)
}
