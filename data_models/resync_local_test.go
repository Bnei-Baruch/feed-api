//go:build localnats

package data_models

import (
	"database/sql"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/events"
)

// Manual test: local nats streaming container "stan-auto-test" (memory store) on 127.0.0.1:24222
// and local postgres on 127.0.0.1:25432 with mdb_auto_survive_nats / localmdb_auto_survive_nats.
func TestResyncRunsAfterNatsReconnect(t *testing.T) {
	open := func(name string) *sql.DB {
		db, err := sql.Open("postgres", "postgres://postgres:test@127.0.0.1:25432/"+name+"?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
		return db
	}
	remote, local := open("mdb_auto_survive_nats"), open("localmdb_auto_survive_nats")

	viper.Set("nats.url", "nats://127.0.0.1:24222")
	viper.Set("nats.client-id", "feed-api-test")
	viper.Set("nats.cluster-id", "test-cluster")
	viper.Set("nats.subject", "mdb2")
	viper.Set("nats.durable", true)
	viper.Set("nats.durable-name", "feed-api-events")
	defer events.RunListener()()

	mv := MakeMdbView(local, remote)
	hook := test.NewGlobal()
	if err := mv.Refresh(); err != nil || logged(hook, "syncing local mdb") {
		t.Fatalf("refresh before restart: err=%v, resync ran=%v", err, logged(hook, "syncing local mdb"))
	}

	if out, err := exec.Command("docker", "restart", "stan-auto-test").CombinedOutput(); err != nil {
		t.Fatalf("restart: %v %s", err, out)
	}
	for end := time.Now().Add(150 * time.Second); !logged(hook, "Finished syncing all tables."); time.Sleep(time.Second) {
		if time.Now().After(end) {
			t.Fatal("resync did not run after nats reconnect")
		}
		if err := mv.Refresh(); err != nil {
			t.Fatalf("refresh after restart: %v", err)
		}
	}
	t.Log("resync ran after reconnect: ", logged(hook, "Reconnected to nats, syncing local mdb."))
}

func logged(hook *test.Hook, msg string) bool {
	for _, e := range hook.AllEntries() {
		if strings.Contains(e.Message, msg) && e.Level <= log.InfoLevel {
			return true
		}
	}
	return false
}
