package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"syscall"
)

func main() {
	log.Println("Starting...")

	log.Printf("\n%s\n", read("/vault/secrets/token"))
	log.Printf("\n%s\n", read("/home/nonroot/.vault-token"))
	log.Printf("\n%s\n", read("/vault/secrets/config.txt"))

	user, err := user.Current()
	if err != nil {
		log.Printf("error getting current user: %v", err)
	} else {
		log.Printf("username: %v, uid: %v, gid: %v", user.Username, user.Uid, user.Gid)
	}
	log.Println("Request: :8080/<name>")
	log.Println("Listening on :8080...")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Request: %s\n", r.URL.Path[:])
	fmt.Fprintf(w, "%s\n", read(r.URL.Path[:]))
}

func read(path string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Start read: %s\n", path))
	info, err := os.Stat(path)
	if err != nil {
		sb.WriteString(fmt.Sprintf("error getting stats on file %s: %v", path, err))
	} else {
		stat := info.Sys().(*syscall.Stat_t)
		sb.WriteString(fmt.Sprintf("uid: %v, gid: %v, mode: %o\n", stat.Uid, stat.Gid, stat.Mode))

		content, err := os.ReadFile(path)
		if err != nil {
			sb.WriteString(fmt.Sprintf("error reading file %s: %v", path, err))
		} else {
			sb.WriteString(fmt.Sprintf("contents: %s", path))
			sb.WriteString(fmt.Sprintln(string(content)))
			sb.WriteString(fmt.Sprintln("Stop read..."))
		}
	}
	return sb.String()
}
