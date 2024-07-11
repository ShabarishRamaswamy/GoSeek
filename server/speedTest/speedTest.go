package speedTest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type NetworkSpeed struct {
	Time float64 `json:"time"`
}

func SpeedTest(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	if category == "request" {
		err := Test_Client_Speed(w, r)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else if category == "response" && r.Method == "POST" {
		networkSpeedClient, err := Get_Client_Speed(r.Body)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
		}
		finalNetworkSpeedInMBs := 1000 / networkSpeedClient.Time
		fmt.Printf("%+v, %+v MB/s", networkSpeedClient, finalNetworkSpeedInMBs)
	}
}

// This will be converted into a separate repository in a bit.
// For now, we will assume that speed of the client DOES NOT VARY.
func Test_Client_Speed(w http.ResponseWriter, r *http.Request) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while reading file: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	contents, err := os.ReadFile(filepath.Join(wd, "assets", "1MB_Test.txt"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Write(contents)
	return nil
}

func Get_Client_Speed(resBody io.ReadCloser) (NetworkSpeed, error) {
	resp, err := io.ReadAll(resBody)
	if err != nil {
		return NetworkSpeed{}, err
	}

	fmt.Println(string(resp))

	var ns NetworkSpeed
	err = json.Unmarshal(resp, &ns)
	if err != nil {
		return NetworkSpeed{}, err
	}
	return ns, nil
}
