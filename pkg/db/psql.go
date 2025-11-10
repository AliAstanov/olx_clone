package db

import (
	"context"
	"fmt"
	"os"

	"github.com/AliAstanov/olx_clone/config"
	"github.com/jackc/pgx/v5"
)

func ConnToDb(pgCfg config.PgConfig) (*pgx.Conn, error) {

	ctx := context.Background()

	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		pgCfg.Username,
		pgCfg.Password,
		pgCfg.Host,
		pgCfg.Port,
		pgCfg.DatabaseName,
	)

	db, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v \n", err)
		return nil, err
	}
	return db, nil

}
