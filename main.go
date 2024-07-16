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

func main() {
	wd, _ := os.Getwd()
	ctx := context.Background()
	fmt.Println("Serving Static files from: ", filepath.Join(wd, "assets"))

	httpWebserver := structs.GetHTTPWebserver(ctx, wd)
	routers := router.GetNewRouter(*httpWebserver)

	r := routers.InitializeAllRoutes()
	http.ListenAndServe("localhost:8000", r)
}
