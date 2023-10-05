package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	tokens map[string]bool
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TokenGen V-2.0")
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Generar un token aleatorio
	token := randomToken()

	// Almacenar el token en el mapa de tokens
	tokens[token] = true

	// Devolver el token en la respuesta
	fmt.Fprintf(w, "%s", token)
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	mutex.Lock()
	defer mutex.Unlock()

	// Validar si el token existe en el mapa de tokens
	if tokens[token] {
		// Eliminar el token de la memoria
		delete(tokens, token)
		fmt.Fprint(w, "true")
	} else {
		fmt.Fprint(w, "false")
	}
}

func randomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenLength := 32
	b := make([]byte, tokenLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	tokens = make(map[string]bool)

	http.HandleFunc("/", handler)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/validate-token", validateTokenHandler)
	http.ListenAndServe(":5000", nil)
}
