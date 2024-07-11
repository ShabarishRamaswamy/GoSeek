package routers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	custom "github.com/ShabarishRamaswamy/GoSeek/server/customDefault"
	"github.com/ShabarishRamaswamy/GoSeek/server/speedTest"
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
	r.HandleFunc("/speedTest/{category}", speedTest.SpeedTest)
	r.HandleFunc("/default", defaultImplementation)
	r.HandleFunc("/custom", customImeplementation)
	r.PathPrefix("/serve/").HandlerFunc(custom.ServeCustom)

	// r.Use(utils.LoggingMiddleware)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	return r
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}

	indexFilePath := filepath.Join(wd, "frontend", "index.html")
	template.Must(template.ParseFiles(indexFilePath)).Execute(w, vp)
}

func defaultImplementation(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	wd, _ := os.Getwd()

	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	defaultImplementationPath := filepath.Join(wd, "frontend", "default", "default.html")

	template.Must(template.ParseFiles(defaultImplementationPath)).Execute(w, vp)
}

func customImeplementation(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()

	vp := VideoPath{VideoPath: "/serve/BBB-Test-Video", ImgPath: "/assets/linux-test-img.png"}
	customImplementationPath := filepath.Join(wd, "frontend", "custom", "custom.html")

	template.Must(template.ParseFiles(customImplementationPath)).Execute(w, vp)
}
