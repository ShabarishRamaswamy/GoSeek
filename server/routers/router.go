package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	speedtest "github.com/ShabarishRamaswamy/GoSeek/server/speedTest"
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
	r.HandleFunc("/speedTest/{category}", speedTest)
	r.HandleFunc("/default", defaultImplementation)
	r.HandleFunc("/custom", customImeplementation)
	r.PathPrefix("/serve/").HandlerFunc(serveCustom)

	// r.Use(LoggingMiddleware)
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

// TODOS;
// Errors when reading files.
// Errors when parsing ranges. [File Sizes, If Ranges are not correctly formatted, If no range is present].
func serveCustom(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	fmt.Printf("Got the Request %+v\n", r)
	filePath := filepath.Join(wd, "assets", "BBB-Test-Video.mp4")

	// fmt.Printf("\nRequest headers: %+v\n", r.Header)

	videoFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	defer videoFile.Close()

	videoFileProperties, err := videoFile.Stat()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}

	// fmt.Printf("\nRange headers: %+v\n", r.Header["Range"])
	requestRangeUnprocessed := r.Header["Range"]
	var reqRange []string
	var vidSeekStart, vidSeekEnd int

	if requestRangeUnprocessed[0] != "" {
		reqRange = strings.Split(requestRangeUnprocessed[0], "=")
		reqRange = strings.Split(reqRange[1], "-")
		// fmt.Println("Request Range: ", reqRange)

		vidSeekStart, _ = strconv.Atoi(reqRange[0])
		if reqRange[0] == "0" || reqRange[1] == "" {
			// fmt.Println("Initial Request")
			vidSeekEnd = int(videoFileProperties.Size()) - 1
		} else {
			vidSeekEnd, _ = strconv.Atoi(reqRange[1])
		}
	}

	// fmt.Println("Seek Properties: ", vidSeekStart, vidSeekEnd)

	vidSeekSize := vidSeekEnd - vidSeekStart + 1

	w.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", vidSeekStart, vidSeekEnd, videoFileProperties.Size()))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", vidSeekSize))
	w.Header().Add("Content-Type", "video/mp4")
	w.Header().Add("Last-Modified", videoFileProperties.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Add("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	videoFile.Seek(int64(vidSeekStart), 0)
	io.CopyN(w, videoFile, int64(vidSeekSize))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		fmt.Printf("\n\nURL: %s\nReq: %+v\n\n", r.RequestURI, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func speedTest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n%+v, %+v\n\n", r.RequestURI, mux.Vars(r))

	category := mux.Vars(r)["category"]
	if category == "request" {
		err := speedtest.Test_Client_Speed(w, r)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else if category == "response" && r.Method == "POST" {
		networkSpeedClient, err := speedtest.Get_Client_Speed(r.Body)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
		}
		finalNetworkSpeedInMBs := 1000 / networkSpeedClient.Time
		fmt.Printf("%+v, %+v MB/s", networkSpeedClient, finalNetworkSpeedInMBs)
	}
}
