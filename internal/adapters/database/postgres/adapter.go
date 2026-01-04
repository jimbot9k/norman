package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PostgresAdapter struct {
	conn *pgx.Conn
}

func (a *PostgresAdapter) Name() string {
	return "PostgreSQL"
}

func (a *PostgresAdapter) Version() string {
	return "v1"
}

func (a *PostgresAdapter) UniqueSignature() string {
	return a.Name() + "-" + a.Version()
}

func (a *PostgresAdapter) IsConnectionStringCompatible(connString string) bool {
	if strings.Contains(connString, "@tcp(") {
		return false
	}

	_, err := pgx.ParseConfig(connString)
	if err != nil {
		return false
	}

	return true
}

func (a *PostgresAdapter) Connect(connString string) error {
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return err
	}
	a.conn = conn
	return nil
}

func (a *PostgresAdapter) Close() error {
	return a.conn.Close(context.Background())
}

func (a *PostgresAdapter) IsConnected() bool {
	return a.conn != nil
}
