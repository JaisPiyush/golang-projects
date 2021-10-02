package pkg

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

type UrlVariable = map[string]interface{}

type Params struct {
	Variables UrlVariable
	Query     QueryDict
	Req       *http.Request
	Res       http.ResponseWriter
}

type HandleFunction = func(*Params)

type RoutePattern struct {
	Pattern map[string]string
	Handler HandleFunction
	Method  string
}
type RouteInterface interface {
	extractVariables(url string, dummyIndex []int, req *http.Request) (*UrlVariable, error)
	ExtractVariables(url string, dummyUrl string, req *http.Request) (*UrlVariable, error)
	Get(pattern string, handler HandleFunction)
	Post(pattern string, handler HandleFunction)
}

type Router struct {
	routes                map[string][]*RoutePattern
	MethodNotFoundHandler http.HandlerFunc
	ErrorHandler          http.HandlerFunc
}

type RouterInterface interface {
	GetHandler(url string, req *http.Request) (*UrlVariable, error)
	ExecuteRequest(res http.ResponseWriter, req *http.Request)
}

func NewRouter() *Router {
	return &Router{routes: make(map[string][]*RoutePattern)}
}

func parseDynamicPattern(pattern string) *map[string]string {
	dynamicParts := make(map[string]string)
	for _, part := range strings.Split(pattern, "/") {
		if string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
			subParts := strings.Split(part, ":")
			if len(subParts) == 2 {
				dynamicParts[subParts[0]] = subParts[1]
			} else {
				dynamicParts[subParts[0]] = ""
			}
		}
	}
	return &dynamicParts
}

func (router *Router) Get(pattern string, handler HandleFunction) {
	dynamicParts := parseDynamicPattern(pattern)
	dynamicUrl := GetDynamicUrl(pattern)
	if _, err := router.routes[dynamicUrl]; !err {
		router.routes[dynamicUrl] = make([]*RoutePattern, 0)
	}
	router.routes[dynamicUrl] = append(router.routes[dynamicUrl], &RoutePattern{
		Pattern: *dynamicParts,
		Handler: handler,
		Method:  "GET",
	})
}

func (router *Router) Post(pattern string, handler HandleFunction) {
	dynamicParts := parseDynamicPattern(pattern)
	dynamicUrl := GetDynamicUrl(pattern)
	if _, err := router.routes[dynamicUrl]; !err {
		router.routes[dynamicUrl] = make([]*RoutePattern, 0)
	}
	router.routes[dynamicUrl] = append(router.routes[dynamicUrl], &RoutePattern{
		Pattern: *dynamicParts,
		Handler: handler,
		Method:  "POST",
	})
}

// Search through the route pattern table and returns

func (router Router) GetHandler(url string, req *http.Request) (*UrlVariable, HandleFunction, error) {
	dummyUrl := GetDynamicUrl(url)
	for key, routes := range router.routes {
		if dummyUrl == key {
			for index := range routes {
				if urlVariable, err := router.routes[key][index].ExtractVariables(url, dummyUrl, req); err == nil {
					return urlVariable, router.routes[key][index].Handler, nil
				}
			}
			return &UrlVariable{}, nil, errors.New("no method found")
		}
	}
	return &UrlVariable{}, nil, errors.New("no method found")
}

func (router Router) ExecuteRequest(res http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	defer func() {
		if r := recover(); r != nil {
			router.ErrorHandler(res, req)
		}
	}()
	if urlVariables, handler, err := router.GetHandler(url, req); err == nil {
		handler(&Params{Req: req, Res: res, Variables: *urlVariables, Query: *QueryDictFromRawQuery(req.URL.RawQuery)})
	} else {
		router.MethodNotFoundHandler(res, req)
	}

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

func (route RoutePattern) ExtractVariables(url string, dummyUrl string, req *http.Request) (*UrlVariable, error) {
	dummyIndex := make([]int, 0)
	for index, str := range strings.Split(dummyUrl, "/") {
		if str == "__dynamic__" {
			dummyIndex = append(dummyIndex, index)
		}
	}
	return route.extractVariables(url, dummyIndex, req)
}

func (route RoutePattern) extractVariables(url string, dummyIndex []int, req *http.Request) (*UrlVariable, error) {
	urlVariables := UrlVariable{}
	if len(dummyIndex) != len(route.Pattern) || route.Method != req.Method {
		return &urlVariables, errors.New("pattern or method is not matching")
	}

	splitUrl := strings.Split(url, "/")
	keyIndex := 0
	for key := range route.Pattern {
		if len(route.Pattern[key]) == 0 {
			urlVariables[key] = splitUrl[dummyIndex[keyIndex]]
		} else {
			re := regexp.MustCompile(route.Pattern[key])
			value := re.FindString(splitUrl[dummyIndex[keyIndex]])
			if len(value) > 0 && value == splitUrl[dummyIndex[keyIndex]] {
				urlVariables[key] = value
			} else {
				return &urlVariables, errors.New("pattern is not matching")
			}
		}
		keyIndex += 1
	}
	return &urlVariables, nil

}
