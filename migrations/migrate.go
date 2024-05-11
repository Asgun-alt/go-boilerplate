package migrations

import (
	"fmt"
	"go-boilerplate/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func Migrate(args []string) {
	if len(args) < 1 {
		panic("missing argument: ./{bin-file} [goose-command]")
	}

	cfg := config.New()
	if cfg == nil {
		panic("cannot load config")
	}

	dbConnection := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
		cfg.Infrastructure.Database.Dialect,
		cfg.Infrastructure.Database.Username,
		cfg.Infrastructure.Database.Password,
		cfg.Infrastructure.Database.Host,
		cfg.Infrastructure.Database.DBName,
	)

	// reading all custom-defined args
	migrationDir, gooseCommand := "migrations/sql", args[0]

	// check db connection
	db, err := goose.OpenDBWithDriver(cfg.Infrastructure.Database.Dialect, dbConnection)
	if err != nil {
		panic(fmt.Sprintf("goose: failed to open DB: %v\n", err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(fmt.Sprintf("goose: failed to close DB: %v\n", err))
		}
	}()

	// executing actual goose
	if errGoose := goose.Run(gooseCommand, db, migrationDir, args...); errGoose != nil {
		panic(fmt.Sprintf("goose %v: %v", gooseCommand, errGoose))
	}
}
