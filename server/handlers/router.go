package router

import (
	"net/http"
	"path/filepath"
	"text/template"

	custom "github.com/ShabarishRamaswamy/GoSeek/server/customDefault"
	"github.com/ShabarishRamaswamy/GoSeek/server/speedTest"
	"github.com/ShabarishRamaswamy/GoSeek/structs"
	"github.com/gorilla/mux"
)

type VideoPath struct {
	VideoPath string
	ImgPath   string
}

type Router struct {
	Webserver structs.HTTPWebserver
}

func GetNewRouter(ws structs.HTTPWebserver) *Router {
	return &Router{
		Webserver: ws,
	}
}

func (router Router) InitializeAllRoutes() *mux.Router {
	// fmt.Println("Initializing Routers")
	r := mux.NewRouter()
	r.HandleFunc("/", router.indexPage)
	r.HandleFunc("/speedTest/{category}", speedTest.SpeedTest)
	r.HandleFunc("/default", router.defaultImplementation)
	r.HandleFunc("/custom", router.customImeplementation)
	r.HandleFunc("/http-live-streaming", router.hls)
	r.HandleFunc("/dash", router.dash)
	r.PathPrefix("/serve/").HandlerFunc(custom.ServeCustomHandler)

	// r.Use(utils.LoggingMiddleware)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	return r
}

func (router Router) indexPage(w http.ResponseWriter, r *http.Request) {
	indexFilePath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "index.html")
	template.Must(template.ParseFiles(indexFilePath)).Execute(w, nil)
}

func (router Router) defaultImplementation(w http.ResponseWriter, r *http.Request) {
	vp := VideoPath{VideoPath: "/assets/BBB-Test-Video.mp4", ImgPath: "/assets/linux-test-img.png"}
	defaultImplementationPath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "default", "default.html")

	template.Must(template.ParseFiles(defaultImplementationPath)).Execute(w, vp)
}

func (router Router) customImeplementation(w http.ResponseWriter, r *http.Request) {
	vp := VideoPath{VideoPath: "/serve/BBB-Test-Video", ImgPath: "/assets/linux-test-img.png"}
	customImplementationPath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "custom", "custom.html")

	template.Must(template.ParseFiles(customImplementationPath)).Execute(w, vp)
}

func (router Router) hls(w http.ResponseWriter, r *http.Request) {
	vp := VideoPath{VideoPath: "/assets/HLS_Video/BBB.m3u8", ImgPath: "/assets/linux-test-img.png"}
	customImplementationPath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "hls", "hls.html")

	template.Must(template.ParseFiles(customImplementationPath)).Execute(w, vp)
}

func (router Router) dash(w http.ResponseWriter, r *http.Request) {
	vp := VideoPath{VideoPath: "/assets/DASH_Video/manifest.mpd"}
	customImplementationPath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "dash", "dash.html")

	template.Must(template.ParseFiles(customImplementationPath)).Execute(w, vp)
}
