package tests

import (
	"testing"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func TestRouter(t *testing.T) {
	router := pkg.NewRouter()
	router.Get("/", func(hp *pkg.HttpParams) {
	})
}
