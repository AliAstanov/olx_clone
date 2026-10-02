package db

import (
	"context"
	"fmt"
	"os"

	"github.com/AliAstanov/olx_clone/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnToDb(pgCfg config.PgConfig) (*pgxpool.Pool, error) {

	ctx := context.Background()

	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		pgCfg.Username,
		pgCfg.Password,
		pgCfg.Host,
		pgCfg.Port,
		pgCfg.DatabaseName,
	)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v \n", err)
		return nil, err
	}
	return db, nil

}
