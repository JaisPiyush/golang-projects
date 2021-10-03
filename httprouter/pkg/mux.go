package pkg

import (
	"net/http"
)

type Router struct {
	route                 *Route
	ErrorHandler          http.HandlerFunc
	NotMethodFoundHandler http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		route: &Route{
			Name:          "",
			Regex:         "",
			IsDynamic:     false,
			Children:      make([]*Route, 0),
			MethodHandler: make(map[string]HandleFunction),
		},
	}
}

type RouterInterface interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request)
	Get(url string, handler HandleFunction)
	Post(url string, handler HandleFunction)
	ExecuteRequest(res http.ResponseWriter, req *http.Request)
	SprintRoutes() string
}

func (router Router) SprintRoutes() string {
	return router.route.String()
}

func (router *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	router.ExecuteRequest(res, req)
}

func (router *Router) Get(url string, handler HandleFunction) {
	router.route.AddNewPattern(url, "GET", handler)
}

func (router *Router) Post(url string, handler HandleFunction) {
	router.route.AddNewPattern(url, "POST", handler)
}

func (router Router) ExecuteRequest(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			router.ErrorHandler(res, req)
		}
	}()

	if err := router.route.ExecuteHandler(&HttpParams{
		Response: res,
		Request:  req,
		Url:      req.URL,
		Method:   req.Method,
		Args:     map[string]string{},
	}); err != nil {
		router.NotMethodFoundHandler(res, req)
	}
}
