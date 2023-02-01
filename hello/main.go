package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	// readSecret("hello/config")
	log.Println("Request: :8080/<name>")
	log.Println("Listening on :8080...")

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Request: %s\n", r.URL.Path[1:])
	fmt.Fprintf(w, "%s", readSecret(r.URL.Path[1:]))
}

func readSecret(path string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln("Start Vault.."))
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		sb.WriteString(fmt.Sprintf("unable to initialize Vault client: %v", err))
		return sb.String()
	}

	token, err := readToken("/vault/secrets/token")
	if err != nil {
		sb.WriteString(fmt.Sprintf("unable to read token: %v", err))
		return sb.String()
	}
	client.SetToken(token)

	// Read a secret
	secret, err := client.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		sb.WriteString(fmt.Sprintf("unable to read secret: %v", err))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("secret.Data: %v\n", secret.Data))

	sb.WriteString(fmt.Sprintln("Stop Vault.."))
	return sb.String()
}

func readToken(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
