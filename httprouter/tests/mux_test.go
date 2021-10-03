package tests

import (
	"testing"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func TestRouter(t *testing.T) {
	router := pkg.NewRouter()
	router.Get("/", func(hp *pkg.HttpParams) {

	})

	// router.Get("/user", func(hp *pkg.HttpParams) {})
	router.Get("/user/{key}", func(hp *pkg.HttpParams) {})
	router.Post("/user/{key}", func(hp *pkg.HttpParams) {})
	router.Get("/name", func(hp *pkg.HttpParams) {})
	router.Get("/user/{key}/{id}", func(hp *pkg.HttpParams) {})

	t.Error(router.SprintRoutes())
}
