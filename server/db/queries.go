package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/ShabarishRamaswamy/GoSeek/structs"
	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME string = "db"
const CREATE_DB string = `CREATE TABLE IF NOT EXISTS users (username text not null, email text not null primary key, password text not null);`

func Setup(wd string) *sql.DB {
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
	if username == "" || email == "" || password == "" {
		return errors.New("inputs cannot be empty")
	}

	saveQ := `INSERT INTO users(username, email, password) VALUES(?, ?, ?)`
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

func FindUser(db *sql.DB, email string) (structs.User, error) {
	if email == "" {
		return structs.User{}, errors.New("cannot have empty fields")
	}

	findQ := "SELECT * from users WHERE email = ?"
	statement, err := db.Prepare(findQ)
	if err != nil {
		return structs.User{}, err
	}
	defer statement.Close()

	var usernameDB, emailDB, passwordDB string
	err = statement.QueryRow(email).Scan(&usernameDB, &emailDB, &passwordDB)
	if err != nil {
		return structs.User{}, err
	}
	var User structs.User
	User.Email = emailDB
	User.Password = passwordDB
	User.Name = usernameDB

	// fmt.Println(emailDB, usernameDB, passwordDB)
	return User, nil
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
