package main

import "bytes"
import "encoding/json"
import "fmt"
import "log"
import "net/http"
import "os"
import "os/exec"

type TokenReply struct {
	Status bool
	Token string
	Error string
}

func main() {
	listener, listener_ok := os.LookupEnv("LISTENER")
	if (!listener_ok) {
		listener = ":8080"
	}
	worker_key, worker_key_ok := os.LookupEnv("WORKER_KEY")
	if (!worker_key_ok) {
		fmt.Printf("Warning: WORKER_KEY is not defined\n")
	}
	manager_key, manager_key_ok := os.LookupEnv("MANAGER_KEY")
	if (!manager_key_ok) {
		fmt.Printf("Warning: MANAGER_KEY is not defined\n")
	}

	http.HandleFunc("/worker", func(w http.ResponseWriter, r *http.Request) {
		handleTokenReq(w, r, "worker", worker_key, worker_key_ok)
	})
	http.HandleFunc("/manager", func(w http.ResponseWriter, r *http.Request) {
		handleTokenReq(w, r, "manager", manager_key, manager_key_ok)
	})

	fmt.Printf("Starting server on %s.\n", listener)
	log.Fatal(http.ListenAndServe(listener, nil))
}

func getToken(token_type string) ([]byte, error) {
	token, err := exec.Command("docker", "swarm", "join-token", token_type, "-q").Output()
	if (err != nil) {
		return nil, err
	}
	token = bytes.TrimSuffix(token, []byte("\n"))
	return token, nil
}

func handleTokenReq(w http.ResponseWriter, r *http.Request, token_type string, key string, key_ok bool) () {
	w.Header().Set("Content-Type", "application/json")

	header_key := r.Header.Get("X-Key")
	fmt.Printf("Processing request: From %s, Type %s, Key %s\n", r.RemoteAddr, token_type, header_key)

	if (!key_ok) {
		jData, err := json.Marshal(TokenReply{false, "", "Server not configured for token type"})
		if err != nil {
			panic(err)
			return
		}
		w.Write(jData)
		fmt.Printf("Warning: Server not configured for this token type\n")
		return
	}
	if (header_key == key) {
		// get and return the token
		token, err := getToken(token_type)
		if (err != nil) {
			log.Fatal(err)
		}
		jData, err := json.Marshal(TokenReply{true, string(token), ""})
		if err != nil {
			panic(err)
			return
		}
		w.Write(jData)

	} else {
		// invalid key in request
		jData, err := json.Marshal(TokenReply{false, "", "Invalid token"})
		if err != nil {
			panic(err)
			return
		}
		w.Write(jData)
		fmt.Printf("Warning: invalid key\n")
	}
}

