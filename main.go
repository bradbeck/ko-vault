package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	log.Default().Println("Starting...")
	vaultInteraction()
	log.Println("Request: :8080/<name>")
	log.Println("Listening on :8080...")

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func vaultInteraction() {
	log.Println("Start Vault..")
	config := vault.DefaultConfig()
	config.Address = "http://external-vault:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken("root")
	secretData := map[string]interface{}{
		"password": "Hashi123",
	}

	ctx := context.Background()

	// Write a secret
	_, err = client.KVv2("secret").Put(ctx, "my-secret-password", secretData)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}

	log.Println("Secret written successfully.")

	// Read a secret
	secret, err := client.KVv2("secret").Get(ctx, "my-secret-password")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	value, ok := secret.Data["password"].(string)
	if !ok {
		log.Fatalf("value type assertion failed: %T %#v", secret.Data["password"], secret.Data["password"])
	}

	if value != "Hashi123" {
		log.Fatalf("unexpected password value %q retrieved from vault", value)
	}

	log.Println("Access granted!")
	log.Println("Stop Vault..")
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Request: %s\n", r.URL.Path[1:])
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
