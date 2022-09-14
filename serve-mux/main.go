package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	//handler implementujacy http.Handler z paczki http
	rh := http.RedirectHandler("http://example.org", 307)

	//handler zaimplenentowany jako struktura implementujaca interfajs http.Handler
	th := timeHandler{format: time.RFC1123}

	//sposob na stworzenie handlera ze zwyklej funkcji
	th2 := http.HandlerFunc(timeHandler2)

	mux.Handle("/foo", rh)
	mux.Handle("/time", th)
	mux.Handle("/time2", th2)
	//skrot na stowrzenie handlera ze zwyklej funkcji i przekazanie go do routera
	mux.HandleFunc("/time3", timeHandler2)
	mux.Handle("/time4", timeHandler3(time.RFC1123))
	mux.Handle("/time5", timeHandler4(time.RFC1123))
	mux.Handle("/time6", timeHandler5(time.RFC1123))

	//
	log.Print("Listening...")

	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	http.ListenAndServe(":3000", mux)
}

type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

// nie trzeba robic specjalnie typu zeby implementowac ServeHTTP
// mozna zrobic dowolna przymujaca i zwracajaca to co ServerHTTP i przekazac to do
// http.HandlerFunc()
func timeHandler2(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}
func timeHandler3(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}
func timeHandler4(format string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	})
}
func timeHandler5(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
}
