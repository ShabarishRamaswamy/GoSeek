package router

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	custom "github.com/ShabarishRamaswamy/GoSeek/server/customDefault"
	"github.com/ShabarishRamaswamy/GoSeek/server/db"
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
	r := mux.NewRouter()
	r.HandleFunc("/", router.indexPage)
	r.HandleFunc("/login", router.register)
	r.HandleFunc("/signup", router.register)
	r.HandleFunc("/register", router.register)

	// Streaming
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
	// If not logged in
	indexFilePath := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "index.html")
	template.Must(template.ParseFiles(indexFilePath)).Execute(w, nil)

	// If Logged in:
	// Home Page with uploaded videos.
}

func (router Router) register(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/login" {
		loginT := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "register", "login.html")
		template.Must(template.ParseFiles(loginT)).Execute(w, nil)
	} else if r.RequestURI == "/signup" {
		signupT := filepath.Join(router.Webserver.BaseWorkingDir, "frontend", "register", "signup.html")
		template.Must(template.ParseFiles(signupT)).Execute(w, nil)
	} else if r.RequestURI == "/register" && r.Method == http.MethodPost {
		r.ParseForm()
		fmt.Printf("%+v", r.Form)

		err := db.SaveUser(router.Webserver.DB, r.Form["name"][0], r.Form["email"][0], r.Form["password"][0])
		if err != nil {
			log.Fatalf("Error %s", err.Error())
		}
	}
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
