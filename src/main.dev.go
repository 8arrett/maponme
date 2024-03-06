package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", logHandler(indexRoute))
	mux.HandleFunc("/api/", logHandler(apiRoute))
	mux.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./src/static"))))

	log.Print("Listening on 8080...\n")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Print("logHandler: Dump request failed.")
			return
		}
		log.Println(fmt.Sprintf("%q", d))
		fn(w, r)
	}
}

/**
 *  indexRoute exists for dev env but a reverse proxy will be more performant in prod
 */
func indexRoute(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/favicon.ico" || r.URL.Path == "/robots.txt" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<!DOCTYPE html><head><meta http-equiv=\"REFRESH\" content=\"0; url=/\"></head><body>404</body></html>"))
		return
	}

	switch r.Method {
	case http.MethodGet:

		fileBuf, err := os.ReadFile("./src/static/index.htm")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(fileBuf)

	default:
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
