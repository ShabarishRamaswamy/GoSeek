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
	"github.com/ShabarishRamaswamy/GoSeek/server/utils"
	"github.com/ShabarishRamaswamy/GoSeek/structs"
)

var PORT int = 8000

func main() {
	wd, _ := os.Getwd()
	ctx := context.Background()
	fmt.Println("Serving Static files from: ", filepath.Join(wd, "assets"))

	db := db.Setup(wd)
	if db == nil {
		log.Fatal("Failed to Setup the DB: ", db)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	env, err := utils.SetupEnv(wd)
	if err != nil {
		log.Fatal("Error reading the Env", err)
	}

	httpWebserver := structs.GetHTTPWebserver(ctx, wd, db, env)
	routers := router.GetNewRouter(*httpWebserver)

	fmt.Println("Listening on Port: ", PORT)

	r := routers.InitializeAllRoutes()
	http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), r)
}
