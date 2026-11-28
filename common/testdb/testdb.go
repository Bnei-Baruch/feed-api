// Package testdb builds temporary databases from migrations, for tests (as in Bnei-Baruch/mdb).
package testdb

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Bnei-Baruch/feed-api/databases/mdb/migrations"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type TestDB struct {
	DB        *sql.DB
	serverUrl string
	name      string
}

// Make creates test_<random> on the server of serverUrl (any database, e.g. postgres), applies the migrations
// in dir: golang-migrate files (*.up.sql) or else rambler files (*.sql, "-- rambler up" sections), then the fixtures.
func Make(serverUrl, dir string, fixtures ...string) (_ *TestDB, err error) {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	t := &TestDB{serverUrl: serverUrl, name: "test_" + hex.EncodeToString(b)}
	if err := exec(serverUrl, "CREATE DATABASE "+t.name); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			t.Destroy()
		}
	}()
	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	u.Path = "/" + t.name
	if t.DB, err = sql.Open("postgres", u.String()); err != nil {
		return nil, err
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	rambler := len(files) == 0
	if rambler {
		files, _ = filepath.Glob(filepath.Join(dir, "*.sql")) // Sorted, as rambler applies them.
	}
	for _, f := range append(files, fixtures...) {
		statements := []string(nil)
		if rambler && !utils.StringInSlice(f, fixtures) {
			m, err := migrations.NewMigration(f)
			if err != nil {
				return nil, err
			}
			if statements = m.Up(); len(statements) == 0 {
				return nil, fmt.Errorf("%s: no rambler up section", f)
			}
		} else if b, err := os.ReadFile(f); err != nil {
			return nil, err
		} else {
			statements = []string{string(b)}
		}
		for _, s := range statements {
			if _, err := t.DB.Exec(s); err != nil {
				return nil, fmt.Errorf("%s: %w", f, err)
			}
		}
	}
	return t, nil
}

func (t *TestDB) Destroy() error {
	if t.DB != nil {
		if err := t.DB.Close(); err != nil {
			return err
		}
	}
	return exec(t.serverUrl, "DROP DATABASE "+t.name)
}

func exec(dbUrl, statement string) error {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(statement)
	return err
}
