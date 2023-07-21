package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var dbName string
	var csvName string
	flag.StringVar(&dbName, "db", "", "SQLite database to import to")
	flag.StringVar(&csvName, "csv", "", "CSV file to import from")
	flag.Parse()

	if dbName == "" || csvName == "" {
		flag.PrintDefaults()
		return
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ping failed %s", err)
	}

	stmt, err := db.Prepare("create table if not exists books (id integer primary key autoincrement, title text, author text)")
	if err != nil {
		log.Fatalf("prepare failed %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("exec failed %s", err)
	}

	f, err := os.Open(csvName)
	if err != nil {
		log.Fatalf("error open csv file %s", err)
	}

	r := csv.NewReader(f)
	// Read the header row.

	_, err = r.Read()
	if err != nil {
		log.Fatalf("missing header row(?) %s", err)
	}

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		title := record[1]
		author := record[2]

		stmt, err = db.Prepare("insert into books(title, author) values(?, ?)")
		if err != nil {
			log.Fatalf("insert prepare value %s", err)
		}

		_, err = stmt.Exec(title, author)
		if err != nil {
			log.Fatalf("insert failed(%s): %s", title, err)
		}

	}

}
