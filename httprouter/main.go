package main

import (
	"log"
	"net/http"
	"time"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func main() {
	r := pkg.NewRouter()
	// r.Get("/", func(params *pkg.Params) {
	// 	fmt.Fprintf(params.Res, "Hi\n")
	// })

	// r.Get("/name/{age}", func(p *pkg.Params) {
	// 	fmt.Fprintln(p.Res, p.Variables)
	// })
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
	// log.Println("Listening on Port 8080")
}
