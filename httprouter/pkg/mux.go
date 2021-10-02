package pkg

import (
	"net/http"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

type RouterInterface interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

func (router *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {}
