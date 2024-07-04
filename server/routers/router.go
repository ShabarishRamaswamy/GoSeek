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
	r.HandleFunc("/", indexPage)
	r.HandleFunc("/default", defaultImplementation)
	r.HandleFunc("/custom", customImeplementation)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	return r
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	// fmt.Println("Working Dir is: ", wd)
	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	// vp := VideoPath{VideoPath: "https://www.w3schools.com/html/mov_bbb.mp4"}

	indexFilePath := filepath.Join(wd, "frontend", "index.html")
	template.Must(template.ParseFiles(indexFilePath)).Execute(w, vp)
}

func defaultImplementation(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()

	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	defaultImplementationPath := filepath.Join(wd, "frontend", "default", "default.html")

	template.Must(template.ParseFiles(defaultImplementationPath)).Execute(w, vp)
}

func customImeplementation(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()

	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	customImplementationPath := filepath.Join(wd, "frontend", "custom", "custom.html")

	template.Must(template.ParseFiles(customImplementationPath)).Execute(w, vp)
}
