package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	HttpPort string
)

func main() {
	flag.StringVar(&HttpPort, "http-port", LookupEnvOrString("PORT", HttpPort), "http port to listen to")

	flag.Parse()

	s := NewServer()
	if HttpPort != "" {
		s.HttpPort = HttpPort
	}

	if err := s.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}


const defaultHttpPort = "3000"

type server struct {
	HttpPort string
}

func NewServer() *server {
	return &server{defaultHttpPort}
}

func (s *server) Run() error {
	if s.HttpPort == "" {
		s.HttpPort = defaultHttpPort
	}

	http.HandleFunc("/", s.auth(s.handleIndex()))

	fmt.Printf("Listening on :%s\n", s.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.HttpPort), nil)
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	}
}

func (s *server) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"status": "not-ok"}`))
			return
		}
		h(w, r)
	}
}
