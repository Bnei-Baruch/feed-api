package utils

import (
	"context"
	"database/sql"
	"fmt"
)

// PostgresNoOpenCheck creates new PostgreSQL health check similar to the one provided by health-go
// with a preconfigured, unmanaged sql.DB instance
func PostgresNoOpenCheck(db *sql.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := db.PingContext(ctx)
		if err != nil {
			return fmt.Errorf("PostgreSQL health check failed on ping: %w", err)
		}

		rows, err := db.QueryContext(ctx, "SELECT VERSION()")
		if err != nil {
			return fmt.Errorf("PostgreSQL health check failed on select: %w", err)
		}
		if err = rows.Close(); err != nil {
			return fmt.Errorf("PostgreSQL health check failed on rows closing: %w", err)
		}

		return nil
	}
}
