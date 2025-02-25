package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ShabarishRamaswamy/GoSeek/structs"
)

func SetupEnv(BaseWorkingDir string) (structs.Env, error) {
	envContents, err := os.ReadFile(filepath.Join(BaseWorkingDir, ".env"))
	if err != nil {
		return structs.Env{}, err
	}
	var envStruct structs.Env
	json.Unmarshal(envContents, &envStruct)
	// fmt.Printf("%+v", envStruct)
	if len(envStruct.SALT) == 0 {
		return structs.Env{}, errors.New("SALT cannot be empty")
	}
	return envStruct, nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		fmt.Printf("\n\nURL: %s\nReq: %+v\n\n", r.RequestURI, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func ParseForm(r *http.Request) map[string]string {
	r.ParseForm()
	formContents := map[string]string{}

	for key, value := range r.Form {
		formContents[key] = value[0]
	}
	return formContents
}

func ServeWebpage(path ...string) *template.Template {
	finalPath := filepath.Join(path...)
	return template.Must(template.ParseFiles(finalPath))
}

func GenerateRandomBytes(n uint8) []byte {
	if n == 0 {
		n = 16
	}

	token := make([]byte, n)
	_, err := rand.Reader.Read(token)
	if err != nil {
		return []byte("")
	}
	return token
}
