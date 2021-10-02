package pkg

import (
	"net/http"
	"strings"
)

type UrlVariable = map[string]interface{}

type HandleFunction = func(http.ResponseWriter, *http.Request, *QueryDict)

type RoutePattern struct {
	Pattern map[string]string
	Handler *HandleFunction
}

type Router struct {
	routes                map[string]*RoutePattern
	MethodNotFoundHandler http.Handler
	ErrorHandler          http.Handler
}

type RouteInterface interface {
	GetDynamicUrl(url string) string
	ExtractVariables(url string) (*UrlVariable, error)
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]*RoutePattern)}
}

func GetDynamicUrl(url string) string {
	split_url := strings.Split(url, "/")
	for i, str := range split_url {
		if string(str[0]) == "{" && string(str[len(str)-1]) == "}" {
			split_url[i] = "__dynamic__"
		}
	}
	return strings.Join(split_url, "/")
}

func (route RoutePattern) ExtractVariables(url string) (*UrlVariable, error) {

}
