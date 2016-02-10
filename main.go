package main

import (
	"log"
	"time"

	"github.com/jackc/pgx"
)

func main() {
	log.Println("start sample")

	pgxConfig := pgx.ConnConfig{
		Host:     "localhost",
		Database: "quetest",
		User:     "quetest",
	}
	pgxConnPoolConfig := pgx.ConnPoolConfig{pgxConfig, 3, nil}
	conn, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		log.Fatal(err)
	}

	// declare
	var i int32
	var name string
	var updatedAt time.Time

	// select a row
	err = conn.QueryRow("select id, name from item where id = $1", 2).Scan(&i, &name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[select a row] id: %d name: %s", i, name)

	// select rows
	rows, err := conn.Query("select id, name, updated_at from item")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&i, &name, &updatedAt); err != nil {
			log.Fatal(err)
		}
		log.Printf(
			"[select rows] id: %d name: %s updated: %s",
			i, name, updatedAt.Format("01/02 15:04:05"))
	}

	// update rows in transaction
	tx, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	res, err := tx.Exec("UPDATE item SET updated_at = $1 WHERE id = $2", time.Now(), 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[update a row] updated num rows: %d", res.RowsAffected())

	err = conn.QueryRow("select name, updated_at from item where id = $1", 2).Scan(&name, &updatedAt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[select a row in transaction] name: %s updated: %s", name, updatedAt.Format("01/02 15:04:05"))

	tx.Commit()

	err = conn.QueryRow("select name, updated_at from item where id = $1", 2).Scan(&name, &updatedAt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[select a row] name: %s updated: %s", name, updatedAt.Format("01/02 15:04:05"))
}
