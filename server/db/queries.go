package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME string = "db"
const CREATE_DB string = `CREATE TABLE IF NOT EXISTS users (username text not null primary key, email text, password text);`

func Setup(wd string) *sql.DB {
	// fmt.Println(wd, "server", "db", DB_NAME, ".db", filepath.Join(wd, "server", "db", DB_NAME+".db"))
	db, err := sql.Open("sqlite3", filepath.Join(wd, "server", "db", DB_NAME+".db"))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	_, err = db.Exec(CREATE_DB)
	if err != nil {
		log.Printf("%q: %s\n", err, CREATE_DB)
		return nil
	}
	return db
}

func SaveUser(db *sql.DB, username, email, password string) error {
	saveQ := `insert into users(username, email, password) values(?, ?, ?)`
	_, err := db.Prepare(saveQ)
	if err != nil {
		log.Printf("%q: %s\n", err, saveQ)
		return err
	}
	_, err = db.Exec(saveQ, username, email, password)
	if err != nil {
		log.Printf("%q: %s\n", err, saveQ)
		return err
	}

	fmt.Println("Saved into DB")

	rows, err := db.Query("select * from users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var email string
		var password string
		err = rows.Scan(&username, &email, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(username, email, password)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func RunQueries(db *sql.DB) {
	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
