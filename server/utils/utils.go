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
	"time"

	"github.com/ShabarishRamaswamy/GoSeek/structs"
	"github.com/golang-jwt/jwt"
)

func SetupEnv(BaseWorkingDir string) (structs.Env, error) {
	envContents, err := os.ReadFile(filepath.Join(BaseWorkingDir, ".env"))
	if err != nil {
		return structs.Env{}, err
	}
	var envStruct structs.Env
	json.Unmarshal(envContents, &envStruct)
	if len(envStruct.JWT_Secret_Key) == 0 {
		return structs.Env{}, errors.New("JWT_Secret_Key cannot be empty")
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

// https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60
func CreateJWT(email string, env structs.Env) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(env.JWT_Secret_Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string, env structs.Env) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return env.JWT_Secret_Key, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
