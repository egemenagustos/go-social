package main

import (
	"go-social/internal/db"
	"go-social/internal/env"
	store "go-social/internal/storage"
	"log"
)

func main() {
	addr := env.GetString("DB_MIGRATOR_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")

	conn, err := db.New(
		addr,
		30,
		30,
		"15m",
	)

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
