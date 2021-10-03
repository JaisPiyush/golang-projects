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
	AddChild(part string, method string, handler HandleFunction)

	UpdateRoute(method string, handler HandleFunction)

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
	keys := "["
	for key, _ := range route.MethodHandler {
		keys += key + ", "
	}
	keys += "]"
	childStr += "]"
	return fmt.Sprintln("\nRoute{ Name: " + route.Name + ", Regex: " + route.Regex + ", IsDynamic: " + strconv.FormatBool(route.IsDynamic) +
		",Children: " + childStr + ", Methods: " + keys + " }")
}

func ParsePart(part string, method string, handler HandleFunction) *Route {
	route := Route{Name: "", Regex: "", IsDynamic: false, MethodHandler: make(map[string]HandleFunction)}
	if handler != nil {
		route.MethodHandler[method] = handler
	}
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

func (route *Route) AddChild(part string, method string, handler HandleFunction) *Route {
	newRoute := ParsePart(part, method, nil)
	route.Children = append(route.Children, newRoute)
	return newRoute
}

func (route *Route) UpdateRoute(method string, handle HandleFunction) {
	route.MethodHandler[method] = handle
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

func (route Route) IsDummyRouteMatching(part string) bool {
	return (!route.IsDynamic && route.Name == part) || (part == fmt.Sprintf("{z%s:%s}", route.Name, route.Regex))
}

func (route *Route) AddNewPattern(url, method string, handler HandleFunction) {
	urlPatterns := strings.Split(CleanSplash(url), "/")
	stack := NewRouteStack()
	if route.IsDummyRouteMatching(urlPatterns[0]) {
		if len(urlPatterns) > 1 {
			urlPatterns = urlPatterns[1:]
		} else {
			urlPatterns = make([]string, 0)
		}
		route.addNewPattern(urlPatterns, stack, method, handler)
	}

}

func (route *Route) addNewPattern(parts []string, stack *RouteStack, method string, handler HandleFunction) {
	if len(parts) > 0 && len(route.Children) > 0 {
		prevStackSize := stack.Size()
		for _, child := range route.Children {
			stack.Push(child)
			break
		}
		if prevStackSize == stack.Size() {
			route.AddChild(parts[0], method, nil)
		}
	} else if len(route.Children) == 0 && len(parts) > 0 {
		stack.Push(route.AddChild(parts[0], method, nil))
	} else if len(parts) == 0 {
		route.UpdateRoute(method, handler)
	}

	if !stack.IsEmpty() {
		if len(parts) > 1 {
			parts = parts[1:]
		} else {
			parts = make([]string, 0)
		}
		stack.Pop().addNewPattern(parts, stack, method, handler)
	}
}

func (route *Route) ExecuteHandler(params *HttpParams) error {
	url := params.Url.Path
	urlParts := strings.Split(CleanSplash(url), "/")
	stack := NewRouteStack()
	if len(urlParts) > 1 {
		urlParts = urlParts[1:]
	} else {
		urlParts = make([]string, 0)
	}
	if params, handler, err := route.findRoute(urlParts, params, stack); err == nil {
		handler(params)
		return nil
	}
	return errors.New("no match found")
}

/*
DFS Algorithm to match url

*/
func (route *Route) findRoute(parts []string, params *HttpParams, stack *RouteStack) (*HttpParams, HandleFunction, error) {

	if len(parts) > 0 && len(route.Children) > 0 {
		for _, child := range route.Children {
			if child.IsRouteMatching(parts[0]) {
				stack.Push(child)
			}
		}
	} else if len(parts) > 0 && len(route.Children) == 0 {
		// need to move next node on breadth
		parts = append([]string{parts[0]}, parts...)
	} else if len(parts) == 0 {
		return params, route.MethodHandler[params.Method], nil
	}

	if !stack.IsEmpty() {
		top := stack.Pop()
		if top.IsDynamic && len(parts) > 0 {
			params.Args[top.Name] = parts[0]
		}
		if len(parts) > 1 {
			parts = parts[1:]
		} else {
			parts = make([]string, 0)
		}
		return top.findRoute(parts, params, stack)
	}
	return params, nil, errors.New("no method found")

}
