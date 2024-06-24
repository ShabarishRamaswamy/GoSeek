package routers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
)

type VideoPath struct {
	VideoPath string
	ImgPath   string
}

func InitializeAllRoutes(wd string) *mux.Router {
	fmt.Println("Initializing Routers")
	r := mux.NewRouter()
	r.HandleFunc("/", sayHi)
	// r.HandleFunc("/assets/{path}", printOnly)
	return r
}

func sayHi(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	// fmt.Println("Working Dir is: ", wd)
	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	// vp := VideoPath{VideoPath: "https://www.w3schools.com/html/mov_bbb.mp4"}

	indexFilePath := filepath.Join(wd, "frontend", "index.html")
	hiTemplate := template.Must(template.ParseFiles(indexFilePath))
	hiTemplate.Execute(w, vp)
}

// func printOnly(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Requested: ", r.RequestURI)
// 	w.con
// }
