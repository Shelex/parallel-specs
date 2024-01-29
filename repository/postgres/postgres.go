package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shelex/parallel-specs/repository"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func (pg *Postgres) ShutDown(ctx context.Context) {
	pg.db.Close()
}

func NewPostgresStorage(ctx context.Context, url string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %s", err.Error())
	}

	logrusLogger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.WarnLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	config.ConnConfig.Logger = logrusadapter.NewLogger(logrusLogger)

	log.Println("starting postgres...")

	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	pg := &Postgres{
		ctx: ctx,
		db:  db,
	}

	repository.DB = pg

	return pg, nil
}
