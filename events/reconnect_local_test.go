//go:build localnats

package events

import (
	"os/exec"
	"testing"
	"time"

	"github.com/nats-io/stan.go"
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
	waitFor := func(id string, d time.Duration) {
		for end := time.Now().Add(d); time.Now().Before(end); time.Sleep(200 * time.Millisecond) {
			for _, e := range ReadAndClearEvents() {
				if e.ID == id {
					return
				}
			}
		}
		t.Fatalf("event %s not received", id)
	}

	publish("before")
	waitFor("before", 5*time.Second)
	if TakeResync() {
		t.Fatal("resync set before any reconnect")
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
		got := false
		for end := time.Now().Add(5 * time.Second); time.Now().Before(end) && !got; time.Sleep(200 * time.Millisecond) {
			for _, e := range ReadAndClearEvents() {
				got = got || e.ID == id
			}
		}
		if got {
			t.Logf("delivery resumed after %s, resync=%v", time.Since(start).Round(time.Second), TakeResync())
			return
		}
	}
}
