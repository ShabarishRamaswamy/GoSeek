package custom

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ShabarishRamaswamy/GoSeek/structs"
)

type ServeCustom struct {
	structs.HTTPWebserver
}

// NOTE: This is essentially a mirror of the default golang implementation.

// TODOS;
// Errors when reading files.
// Errors when parsing ranges. [File Sizes, If Ranges are not correctly formatted, If no range is present].
func ServeCustomHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Video File Properties: %+v\n", videoFileProperties)

	fmt.Printf("\nRange headers: %+v\n", r.Header["Range"])
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
	w.Header().Add("Content-Type", "video/ogg")
	w.Header().Add("Last-Modified", videoFileProperties.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Add("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	videoFile.Seek(int64(vidSeekStart), 0)
	io.CopyN(w, videoFile, int64(vidSeekSize))
}
