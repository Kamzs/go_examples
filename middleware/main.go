package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type TestTypeString string

func (t TestTypeString) f() {
	log.Println(t)
}

//cale to HandelFunc to jest asercja typu.
func middleware1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware1 done")
	})
}
func middleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware2 started")
		next.ServeHTTP(w, r)
		fmt.Println("middleware2 done")
	})
}
func middleware3(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware3 started")
		next.ServeHTTP(w, r)
		fmt.Println("middleware3 done")
	})
}
func main() {
	//var t TestTypeString
	t := TestTypeString("its is just type assertion")
	log.Println(t)
	t.f()

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	loggingHandler := newLoggingHandler(logFile)
	//authHandler := httpauth.SimpleBasicAuth("alice", "pa$$word")

	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", loggingHandler((enforceJSONHandler(finalHandler))))

	mux.Handle("/test1", middleware3(middleware2(middleware1())))

	// todo analyse how to create constructor wrapping middlewares
	// todo commit and push code
	//https://www.alexedwards.net/blog/making-and-using-middleware
	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func enforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)
	}
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
