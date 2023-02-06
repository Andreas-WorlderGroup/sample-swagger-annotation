package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

type MySQL interface {
	TableSensor
}

type Dependency struct {
	zerolog.Logger
}

type Configuration struct {
	DSN string
}

type instance struct {
	Dependency
	Configuration

	*sql.DB
}

func New(cfg Configuration, dep Dependency) MySQL {
	dep.Logger.Info().Msg("starting connect to db...")
	conn, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		dep.Logger.Fatal().Msg("failed connect to db")
	}

	if pingErr := conn.Ping(); pingErr != nil {
		dep.Logger.Fatal().Msg("error ping db")
	}

	dep.Logger.Info().Msg("connected to db...")
	return &instance{Configuration: cfg, Dependency: dep, DB: conn}
}

func (x *instance) ExecContext(ctx context.Context, query string, args ...any) (tx *sql.Tx, err error) {
	tx, err = x.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return
}

func CommitTxSql(list_tx []*sql.Tx) {
	for _, tx := range list_tx {
		if tx != nil {
			tx.Commit()
		}
	}
}

func RollbackTxSql(list_tx []*sql.Tx, err error) error {
	for _, tx := range list_tx {
		if tx != nil {
			tx.Rollback()
		}
	}

	return err
}
