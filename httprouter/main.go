package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func main() {
	r := pkg.NewRouter()
	r.NotMethodFoundHandler = func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "No method found for url %s\n ", req.URL.Path)
	}
	r.Get("/", func(params *pkg.HttpParams) {
		fmt.Fprintf(params.Response, "Hi\n")
	})

	r.Get("/name/{age}", func(p *pkg.HttpParams) {
		fmt.Fprintln(p.Response, "I can't believe you are ", p.Args["age"])
	})

	r.Get("/name/{age}/{gender}", func(p *pkg.HttpParams) {
		fmt.Fprintln(p.Response, "You are nice ", p.Args["gender"], " of age ", p.Args["age"])
	})

	r.Get("/name/{id:[0-9]+}/of/{age:[0-9]+}", func(p *pkg.HttpParams) {
		fmt.Fprintln(p.Response, "You are also nice ", p.Args["gender"], " of age ", p.Args["age"])
	})
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
	// log.Println("Listening on Port 8080")
}
