package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	connStr := os.Getenv("POSTGRESS_CONN")
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Invalid Database String Parsed")
		os.Exit(1)
	}
	config.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Error Connecting to the Databse")
	}
	fmt.Println("Connection was successfull")
	return pool
}

func Close(c *pgxpool.Pool) {
	if c != nil {
		c.Close()
	}
}
