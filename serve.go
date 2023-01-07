package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		log.Fatal("STATIC_DIR env not specified (eg: /public/)")
	}

	fs := http.FileServer(http.Dir(staticDir))
	mux := http.NewServeMux()
	mux.Handle("/healthz", health())
	mux.Handle("/", rootPath(fs))

	log.Printf("starting to listen on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != http.ErrServerClosed {
		log.Fatalf("listen error: %+v", err)
	}

	log.Printf("server shutdown successfully")
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}
}

func rootPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
