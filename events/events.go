package events

import (
	"encoding/json"
	"runtime/debug"
	"sync"

	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	eventsMutex sync.Mutex
	events      []Data
	DebugMode   bool
)

func ReadAndClearEvents() []Data {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	ret := make([]Data, len(events))
	copy(ret, events)
	events = []Data(nil)

	return ret
}

func AddEvent(d Data) {
	eventsMutex.Lock()
	defer eventsMutex.Unlock()

	events = append(events, d)
}

// Returns shut down function.
func RunListener() func() {
	log.SetLevel(log.InfoLevel)

	var err error

	log.Info("Initialize connection to nats")
	natsURL := viper.GetString("nats.url")
	natsClientID := viper.GetString("nats.client-id")
	natsClusterID := viper.GetString("nats.cluster-id")
	natsSubject := viper.GetString("nats.subject")
	sc, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	utils.Must(err)

	DebugMode = viper.GetString("nats.mode") == "debug"

	log.Info("Subscribing to nats subject")
	var startOpt stan.SubscriptionOption
	if viper.GetBool("nats.durable") == true {
		startOpt = stan.DurableName(viper.GetString("nats.durable-name"))
	} else {
		startOpt = stan.DeliverAllAvailable()
	}
	_, err = sc.Subscribe(natsSubject, msgHandler, startOpt, stan.SetManualAckMode())
	utils.Must(err)

	return func() { sc.Close() }
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

	// log.Infof("Adding %+v", d)
	AddEvent(d)
}
