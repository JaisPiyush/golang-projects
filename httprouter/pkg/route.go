package pkg

import (
	"net/http"
	"net/url"
	"strings"
)

type Arguments = map[string]string

type HttpParams struct {
	Response http.ResponseWriter
	Request  *http.Request
	Method   string
	Get      *QueryDict
	Post     *QueryDict
	Args     *Arguments
	Url      *url.URL
}

type HandleFunction = func(*HttpParams)

type Route struct {
	Name          string
	Regex         string
	IsDynamic     bool
	Handler       HandleFunction
	Children      []*Route
	Method        string
	PreservedPath string
}
type RouteInterface interface {

	// Function when any new Node is added
	AddChild(part []string, method string, preservedPath string, handler HandleFunction)

	// Traverse through the tree and add new nodes if required
	AddNewPattern(url string, method string, handler HandleFunction)
}

func parsePart(part string, method string, preservedPath string, handler HandleFunction) *Route {
	route := Route{Name: "", Regex: "", IsDynamic: false, Handler: handler, Method: method, PreservedPath: preservedPath + "/" + part}
	route.Children = make([]*Route, 0)
	if len(part) > 0 && string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
		route.IsDynamic = true
		subparts := strings.Split(part, ":")
		route.Name = subparts[0]
		if len(subparts) == 2 {
			route.Regex = subparts[1]
		}
	} else {
		route.Name = part
	}
	return &route
}

func (route *Route) AddChild(part []string, method string, preservedPath string, handler HandleFunction) {
	newRoute := parsePart(part[0], preservedPath, method, nil)
	if len(part) > 1 {
		newRoute.AddChild(part[1:], newRoute.PreservedPath, method, handler)
	} else {
		newRoute.Handler = handler
	}
	route.Children = append(route.Children, newRoute)
}

func (route *Route) AddNewPattern(url string, method string, handler HandleFunction) {

}
