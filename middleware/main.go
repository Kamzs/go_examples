package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
)

type TestTypeString string
type TestAliasString = string

func (t TestTypeString) f() {
	log.Println(t)
}

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)
	}
}

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
func middlware3contructor() func(h http.Handler) http.Handler {
	return middleware3
}
func main() {
	// to nie jest asercja typu.
	// asercja by byla jakbym mial obiekt jakiegos interfejsu i chcial go zrzutowac na typ
	// poAsercjiTypu, czySiePowiodloCzyNie := obiektTypuInterfejsowego.(typ - np. implementujacy interfejs)
	//przyklad ponziej
	var a interface{} = "a"
	b, ok := a.(string)
	if !ok {
		log.Fatal("not succesfull type assertion")
	}
	fmt.Println(b)

	// to nie jest tez alias
	// bo jest tak
	// type TestTypeString string
	// a nie tak
	// type TestTypeString = string
	// alias nie tworzy nowego typu
	var alias TestAliasString = "a"
	fmt.Println(alias)
	fmt.Printf("alias type: %T i.e. not type 'defined' by alias\n", alias)

	t := TestTypeString("its is just type assertion")
	log.Println(t)
	t.f()

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	//middleware3FromConstructor := middlware3contructor()

	loggingHandler := newLoggingHandler(logFile)

	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", loggingHandler(enforceJSONHandler(finalHandler)))

	mux.Handle("/test1", middleware3(middleware2(middleware1())))

	//there is library for wrapping middlewares
	authHandler := httpauth.SimpleBasicAuth("alice", "pa$$word")

	stdChain := alice.New(loggingHandler, authHandler, enforceJSONHandler)

	mux.Handle("/foo", stdChain.Then(finalHandler))
	mux.Handle("/bar", stdChain.Then(finalHandler))

	// todo analyse how to create constructor wrapping middlewares
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

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
