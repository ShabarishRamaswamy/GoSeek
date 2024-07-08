package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ShabarishRamaswamy/GoSeek/server/routers"
)

func main() {
	wd, _ := os.Getwd()
	fmt.Println("Serving Static files from: ", filepath.Join(wd, "assets"))

	r := routers.InitializeAllRoutes(wd)
	http.ListenAndServe("localhost:8000", r)
}
