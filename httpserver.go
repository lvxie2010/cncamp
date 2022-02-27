package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthZ", healthZ)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("start http server failed, error: %s\n", err.Error())
	} else {
		log.Printf("http server is listening on 8080...")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	goRoot := os.Getenv("GOROOT")
	w.Header().Set("GOROOT", goRoot)
	bodyString := ""
	for k, v := range r.Header {
		for _, vv := range v {
			bodyString += fmt.Sprintf("%s: %s\n", k, vv)
			// fmt.Println(bodyString)
			w.Header().Set(k, vv)
		}
	}
	ip := getClinetIP(r)
	log.Printf("Success! client ip address: %s", ip)
	log.Printf("Success! client response code: %d", 200)

	w.Write([]byte(bodyString))
}

func healthZ(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("working..."))
}

func getClinetIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.Split(strings.TrimSpace(forwardedFor), ",")[0]
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
