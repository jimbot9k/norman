package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type MySqlAdapter struct {
	db *sql.DB
}

func (a *MySqlAdapter) Name() string {
	return "MySQL"
}

func (a *MySqlAdapter) Version() string {
	return "v1"
}

func (a *MySqlAdapter) UniqueSignature() string {
	return a.Name() + "-" + a.Version()
}

func (a *MySqlAdapter) IsConnectionStringCompatible(connString string) bool {
	_, err := mysql.ParseDSN(connString)
	return err == nil
}

func (a *MySqlAdapter) Connect(connString string) error {
	cfg, err := mysql.ParseDSN(connString)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	a.db = db
	return nil
}

func (a *MySqlAdapter) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *MySqlAdapter) IsConnected() bool {
	return a.db != nil
}
