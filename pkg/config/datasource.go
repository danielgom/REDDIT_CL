package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type DBConn interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func InitDatabase(c *Config) *pgxpool.Pool {

	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.DB.User,
		c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)

	conn, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalln("could not connect to database,", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalln("could not connect to database,", err.Error())
	}

	return conn
}
