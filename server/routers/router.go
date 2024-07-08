package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

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
	r.PathPrefix("/serve/").HandlerFunc(serveCustom)

	r.Use(loggingMiddleware)
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

func serveCustom(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	fmt.Printf("Got the Request %+v\n", r)
	filePath := filepath.Join(wd, "assets", "BBB-Test-Video.mp4")

	fmt.Printf("\nRequest headers: %+v\n", r.Header)

	videoFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	videoFileProperties, err := videoFile.Stat()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}

	fmt.Printf("\nRange headers: %+v\n", r.Header["Range"])
	requestRangeUnprocessed := r.Header["Range"]
	reqRange := strings.Split(requestRangeUnprocessed[0], "=")
	reqRange = strings.Split(reqRange[1], "-")
	fmt.Println("Request Range: ", reqRange)
	if reqRange[0] == "0" {
		fmt.Println("Initial Request")
	}

	w.Header().Add("Content-Range", fmt.Sprintf("bytes 0-%d/%d", videoFileProperties.Size(), videoFileProperties.Size()))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", videoFileProperties.Size()))
	w.Header().Add("Content-Type", "video/mp4")
	w.Header().Add("Last-Modified", fmt.Sprintf("%+v", time.Now()))
	w.Header().Add("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	var sendContent io.Reader = videoFile
	io.CopyN(w, sendContent, videoFileProperties.Size())
	// fmt.Printf("%+v", w.Header())
}

// func serveCustom(w http.ResponseWriter, r *http.Request) {
// 	wd, _ := os.Getwd()
// 	fmt.Printf("Got the Request %+v\n", r)
// 	filePath := filepath.Join(wd, "assets", "BBB-Test-Video.mp4")

// 	fmt.Printf("\nRequest headers: %+v\n", r.Header)
// 	fmt.Printf("\nRange headers: %+v\n", r.Header["Range"])

// 	videoFile, err := os.Open(filePath)
// 	if err != nil {
// 		http.Error(w, "File not found.", 404)
// 		return
// 	}
// 	defer videoFile.Close()

// 	videoFileProperties, err := videoFile.Stat()
// 	if err != nil {
// 		http.Error(w, "Could not get file properties.", 500)
// 		return
// 	}

// 	fileSize := videoFileProperties.Size()
// 	rangeHeader := r.Header.Get("Range")

// 	if rangeHeader != "" {
// 		// Parse the Range header to determine the start and end byte positions
// 		rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
// 		start, err := strconv.ParseInt(rangeParts[0], 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid range start.", 400)
// 			return
// 		}

// 		var end int64
// 		if len(rangeParts) > 1 && rangeParts[1] != "" {
// 			end, err = strconv.ParseInt(rangeParts[1], 10, 64)
// 			if err != nil {
// 				http.Error(w, "Invalid range end.", 400)
// 				return
// 			}
// 		} else {
// 			end = fileSize - 1
// 		}

// 		if start > end || end >= fileSize {
// 			http.Error(w, "Invalid range.", 416)
// 			return
// 		}

// 		rangeLength := end - start + 1

// 		w.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
// 		w.Header().Add("Content-Length", fmt.Sprintf("%d", rangeLength))
// 		w.Header().Add("Content-Type", "video/mp4")
// 		w.Header().Add("Last-Modified", videoFileProperties.ModTime().UTC().Format(http.TimeFormat))
// 		w.Header().Add("Accept-Ranges", "bytes")
// 		w.WriteHeader(http.StatusPartialContent)

// 		// Serve the requested range
// 		videoFile.Seek(start, 0)
// 		io.CopyN(w, videoFile, rangeLength)
// 	} else {
// 		// Serve the entire file if no range is requested
// 		w.Header().Add("Content-Length", fmt.Sprintf("%d", fileSize))
// 		w.Header().Add("Content-Type", "video/mp4")
// 		w.Header().Add("Last-Modified", videoFileProperties.ModTime().UTC().Format(http.TimeFormat))
// 		w.Header().Add("Accept-Ranges", "bytes")
// 		w.WriteHeader(http.StatusOK)
// 		io.Copy(w, videoFile)
// 	}

// 	fmt.Printf("%+v", w.Header())
// }

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		fmt.Printf("\n\nURL: %s\nReq: %+v\n\n", r.RequestURI, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
