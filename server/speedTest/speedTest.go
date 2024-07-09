package speedtest

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// This will be converted into a separate repository in a bit.
// For now, we will assume that speed of the client DOES NOT VARY.
func Test_Client_Speed(w http.ResponseWriter, r *http.Request) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while reading file: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New(err.Error())
	}

	contents, err := os.ReadFile(filepath.Join(wd, "assets", "1MB_Test.txt"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	fmt.Println(string(contents[:10]))

	w.Write(contents)
	return nil
}
