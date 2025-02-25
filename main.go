package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ShabarishRamaswamy/GoSeek/server/db"
	router "github.com/ShabarishRamaswamy/GoSeek/server/handlers"
	"github.com/ShabarishRamaswamy/GoSeek/structs"
)

var PORT int = 8000

const DB_NAME string = "db"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working dir: ", err.Error())
	}

	ctx := context.Background()
	fmt.Println("Serving Static files from: ", filepath.Join(wd, "assets"))

	db := db.Setup(filepath.Join(wd, "server", "db", DB_NAME+".db"))
	if db == nil {
		log.Fatal("Failed to Setup the DB: ", db)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	httpWebserver := structs.GetHTTPWebserver(ctx, wd, db)
	routers := router.GetNewRouter(*httpWebserver)

	fmt.Println("Listening on Port: ", PORT)

	r := routers.InitializeAllRoutes()
	http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), r)
}
