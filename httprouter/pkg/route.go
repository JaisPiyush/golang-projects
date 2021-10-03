package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Arguments = map[string]string

type HttpParams struct {
	Response http.ResponseWriter
	Request  *http.Request
	Method   string
	Get      *QueryDict
	Post     *QueryDict
	Args     Arguments
	Url      *url.URL
}

type HandleFunction = func(*HttpParams)

type Route struct {
	Name          string
	Regex         string
	IsDynamic     bool
	Children      []*Route
	MethodHandler map[string]HandleFunction
}
type RouteInterface interface {

	// Function when any new Node is added
	AddChild(part []string, method string, handler HandleFunction)

	// Traverse through the tree and add new nodes if required
	AddNewPattern(url string, method string, handler HandleFunction)

	addNewPattern(parts []string, stack *RouteStack, method string, handler HandleFunction)

	// Part is url subpart withour any /
	IsRouteMatching(part string) bool

	CleanSplash(url string) string

	String() string

	ExecuteHandler(params *HttpParams) (*HttpParams, HandleFunction, error)

	findRoute(parts []string, params *HttpParams, stack *RouteStack) (HandleFunction, error)
}

func (route Route) String() string {
	childStr := "["
	if len(route.Children) > 0 {
		for _, child := range route.Children {
			childStr = childStr + ", " + child.String()
		}
	}
	childStr += "]"
	return fmt.Sprintln("Route{ Name: " + route.Name + ", Regex: " + route.Regex + ", IsDynamic: " + strconv.FormatBool(route.IsDynamic) +
		"Children: " + childStr + ", }")
}

func ParsePart(part string, method string, handler HandleFunction) *Route {
	route := Route{Name: "", Regex: "", IsDynamic: false, MethodHandler: make(map[string]func(*HttpParams))}
	route.MethodHandler[method] = handler
	route.Children = make([]*Route, 0)
	if len(part) > 0 && string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
		part = strings.Replace(part, "{", "", 1)
		part = strings.Replace(part, "}", "", 1)
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

func (route *Route) AddChild(part []string, method string, handler HandleFunction) {
	newRoute := ParsePart(part[0], method, nil)
	if len(part) > 1 {
		newRoute.AddChild(part[1:], method, handler)
	} else {
		newRoute.MethodHandler[method] = handler
	}
	route.Children = append(route.Children, newRoute)
}

func CleanSplash(url string) string {
	return strings.TrimSuffix(url, "/")
}

func (route Route) IsRouteMatching(part string) bool {
	if (!route.IsDynamic && route.Name == part) || (route.IsDynamic && route.Regex == "") {
		return true
	} else if route.IsDynamic {
		re := regexp.MustCompile(route.Regex)
		data := re.FindString(part)
		return len(data) == len(part)

	}
	return false
}

func (route *Route) AddNewPattern(url, method string, handler HandleFunction) {
	urlParts := strings.Split(CleanSplash(url), "/")
	stack := NewRouteStack()
	route.addNewPattern(urlParts, stack, method, handler)
}

func (route *Route) addNewPattern(parts []string, stack *RouteStack, method string, handler HandleFunction) {
	if route.IsRouteMatching(parts[0]) {
		// route part is matching and few parts are left
		if len(parts) > 1 {
			// route has not child --> Add child
			if len(route.Children) == 0 {
				route.AddChild(parts[1:], method, handler)
			} else {
				// route has children --> add them to stack, Pop the top out and addNewPattern
				stack.PushArray(route.Children)
				parts = parts[1:]
			}
		}

	}
	// route does not match pop into the next one
	if stack.Size() > 0 {
		stack.Pop().addNewPattern(parts, stack, method, handler)
	}
}

func (route Route) ExecuteHandler(params *HttpParams) error {
	url := params.Url.Path
	urlParts := strings.Split(CleanSplash(url), "/")
	stack := NewRouteStack()
	if params, handler, err := route.findRoute(urlParts, params, stack); err == nil {
		handler(params)
		return nil
	}
	return errors.New("no method handler found")
}

/*
DFS Algorithm to match url

*/
func (route Route) findRoute(parts []string, params *HttpParams, stack *RouteStack) (*HttpParams, HandleFunction, error) {
	if route.IsRouteMatching(parts[0]) {
		if route.IsDynamic {
			params.Args[route.Name] = parts[0]
		}
		if len(parts) > 1 && len(route.Children) > 0 {
			stack.PushArray(route.Children)

		} else if handler, err := route.MethodHandler[params.Method]; len(parts) == 1 && !err {
			if params.Url.RawQuery != "" {
				params.Get = QueryDictFromRawQuery(params.Url.RawQuery)
			}
			// TODO: Need to implement body digester for post requests
			return params, handler, nil
		}
	}
	if stack.Size() > 0 {
		return stack.Pop().findRoute(parts, params, stack)
	}
	return params, nil, errors.New("no method handler found")

}
