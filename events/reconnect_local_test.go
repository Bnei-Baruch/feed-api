//go:build localnats

package events

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
)

// Manual test against a local nats streaming container "stan-auto-test" on 127.0.0.1:24222.
func TestReconnectAfterNatsRestart(t *testing.T) {
	viper.Set("nats.url", "nats://127.0.0.1:24222")
	viper.Set("nats.client-id", "feed-api-test")
	viper.Set("nats.cluster-id", "test-cluster")
	viper.Set("nats.subject", "mdb2")
	viper.Set("nats.durable", true)
	viper.Set("nats.durable-name", "feed-api-events")
	hook := test.NewGlobal()
	shutdown := RunListener()
	defer shutdown()

	publish := func(id string) {
		pc, err := stan.Connect("test-cluster", "publisher-"+id, stan.NatsURL("nats://127.0.0.1:24222"))
		if err != nil {
			t.Fatal(err)
		}
		defer pc.Close()
		if err := pc.Publish("mdb2", []byte(`{"id":"`+id+`","type":"test"}`)); err != nil {
			t.Fatal(err)
		}
	}
	resync := false
	received := func(id string, d time.Duration) bool {
		for end := time.Now().Add(d); time.Now().Before(end); time.Sleep(200 * time.Millisecond) {
			datas, r := ReadAndClearEvents()
			resync = resync || r
			for _, e := range datas {
				if e.ID == id {
					return true
				}
			}
		}
		return false
	}

	publish("before")
	if !received("before", 5*time.Second) || resync {
		t.Fatalf("before restart: event not received or resync=%v", resync)
	}

	if out, err := exec.Command("docker", "restart", "stan-auto-test").CombinedOutput(); err != nil {
		t.Fatalf("restart: %v %s", err, out)
	}
	// Keep publishing until delivery resumes (lost connection is detected within ~70s by pings).
	start := time.Now()
	for i := 0; ; i++ {
		if time.Since(start) > 150*time.Second {
			t.Fatal("delivery did not resume after nats restart")
		}
		id := "after-" + string(rune('a'+i%26))
		if pc, err := stan.Connect("test-cluster", "pub-after", stan.NatsURL("nats://127.0.0.1:24222")); err == nil {
			pc.Publish("mdb2", []byte(`{"id":"`+id+`","type":"test"}`))
			pc.Close()
		}
		if received(id, 5*time.Second) {
			// Resync is flagged just after resubscribing, give it a moment.
			received("", time.Second)
			lost := false
			for _, e := range hook.AllEntries() {
				lost = lost || strings.HasPrefix(e.Message, "Nats connection lost")
			}
			// File store keeps the client registered, so no connection loss and no resync.
			if resync != lost {
				t.Fatalf("resync=%v but connection lost=%v", resync, lost)
			}
			t.Logf("delivery resumed after %s, connection lost=resync=%v", time.Since(start).Round(time.Second), resync)
			return
		}
	}
}
