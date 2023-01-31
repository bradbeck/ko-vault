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
	// config.Address = "http://external-vault:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Printf("unable to initialize Vault client: %v", err)
		return
	}

	client.SetToken("root")
	secretData := map[string]interface{}{
		"username": "appuser",
		"password": "suP3rsec(et!",
	}

	ctx := context.Background()

	// Write a secret
	// _, err = client.KVv2("secret/hello").Put(ctx, "config", secretData)
	// if err != nil {
	// 	log.Printf("unable to write secret: %v", err)
	// 	return
	// }

	log.Println("Secret written successfully.")

	// Read a secret
	secret, err := client.KVv2("secret/hello").Get(ctx, "config")
	if err != nil {
		log.Printf("unable to read secret: %v", err)
		return
	}

	value, ok := secret.Data["password"].(string)
	if !ok {
		log.Printf("value type assertion failed: %T %#v", secret.Data["password"], secret.Data["password"])
		return
	}

	if value != secretData["password"] {
		log.Printf("unexpected %q value %q retrieved from vault", secretData["password"], value)
		return
	}

	log.Println("Access granted!")
	log.Println("Stop Vault..")
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Request: %s\n", r.URL.Path[1:])
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
