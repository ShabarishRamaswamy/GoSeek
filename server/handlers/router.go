package router

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	custom "github.com/ShabarishRamaswamy/GoSeek/server/customDefault"
	"github.com/ShabarishRamaswamy/GoSeek/server/db"
	"github.com/ShabarishRamaswamy/GoSeek/server/speedTest"
	"github.com/ShabarishRamaswamy/GoSeek/server/utils"
	"github.com/ShabarishRamaswamy/GoSeek/structs"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
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
	if r.RequestURI == "/login" && r.Method == http.MethodGet {
		utils.ServeWebpage(router.Webserver.BaseWorkingDir, "frontend", "register", "login.html").Execute(w, nil)
	} else if r.RequestURI == "/signup" && r.Method == http.MethodGet {
		utils.ServeWebpage(router.Webserver.BaseWorkingDir, "frontend", "register", "signup.html").Execute(w, nil)
	} else if r.RequestURI == "/signup" && r.Method == http.MethodPost {
		formContents := utils.ParseForm(r)

		// Ref: https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
		salt := utils.GenerateRandomBytes(16)
		if len(salt) == 0 {
			log.Println("Internal server error: Salt is empty")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hash := argon2.IDKey([]byte(formContents["password"]), salt, 1, 64*1024, 4, 32)

		b64Salt := base64.RawStdEncoding.EncodeToString(salt)
		b64Hash := base64.RawStdEncoding.EncodeToString(hash)

		encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", 1, 64*1024, 4, 32, b64Salt, b64Hash)

		err := db.SaveUser(router.Webserver.DB, formContents["username"], formContents["email"], encodedHash)
		if err != nil {
			log.Printf("Error %s", err.Error())
			w.Write([]byte("Sorry not allowed"))
			return
		}
		log.Printf("User: %s saved successfully\n", formContents["username"])
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.RequestURI == "/login" && r.Method == http.MethodPost {
		formContents := utils.ParseForm(r)

		user, err := db.FindUser(router.Webserver.DB, formContents["email"])
		if err != nil {
			log.Println("Error with login: ", err.Error())
			// TODO: redirect /login?error=invalid_credentials
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		pwd := strings.Split(user.Password, "$")
		saltDB, err := base64.RawStdEncoding.Strict().DecodeString(pwd[len(pwd)-2])
		if err != nil {
			log.Println("User not found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// fmt.Println("Password Form:", formContents["password"])
		hashCur := argon2.IDKey([]byte(formContents["password"]), saltDB, 1, 64*1024, 4, 32)

		hashDB, err := base64.RawStdEncoding.DecodeString(pwd[len(pwd)-1])
		if err != nil {
			log.Println("Internal server error: Salt is empty")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if subtle.ConstantTimeCompare(hashDB, hashCur) != 1 {
			log.Println("Incorrect password. Malicous user detected")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		jwtToken, err := utils.CreateJWT(formContents["email"], router.Webserver.Env)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Auth-Token", jwtToken)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
