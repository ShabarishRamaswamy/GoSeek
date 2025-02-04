package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	router "github.com/ShabarishRamaswamy/GoSeek/server/handlers"
	"github.com/ShabarishRamaswamy/GoSeek/structs"
)

var PORT int = 8000

func main() {
	wd, _ := os.Getwd()
	ctx := context.Background()
	fmt.Println("Serving Static files from: ", filepath.Join(wd, "assets"))

	httpWebserver := structs.GetHTTPWebserver(ctx, wd)
	routers := router.GetNewRouter(*httpWebserver)

	fmt.Println("Listening on Port: ", PORT)

	r := routers.InitializeAllRoutes()
	http.ListenAndServe(fmt.Sprintf("localhost:%d", PORT), r)
}
