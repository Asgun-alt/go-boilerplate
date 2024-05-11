package database

//nolint:revive // just import no need to worry.
import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
}

func New(cfg *Config) Database {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.DBName,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	// err = db.Ping()
	// if err != nil {
	// 	panic(fmt.Sprintf("Cannot ping to database: %v", err))
	// }

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIddleConn)
	db.SetConnMaxLifetime(cfg.MaxLifeTimeConn)

	return Database{
		db,
	}
}
