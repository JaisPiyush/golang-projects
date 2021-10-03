package pkg

import (
	"errors"
	"strconv"
	"strings"
)

type QueryDict struct {
	dict map[string][]interface{}
}

type QueryDictInterface interface {
	// Create() *QueryDict
	Set(key string, values []interface{})
	SetAttr(key string, values ...interface{})
	Get(key string, defaults ...interface{}) []interface{}
	Clone() QueryDict
	Encode() string
}

func NewQueryDict() *QueryDict {
	queryDict := new(QueryDict)
	queryDict.dict = make(map[string][]interface{})
	return queryDict
}

func QueryDictFromRawQuery(query string) *QueryDict {
	fields := strings.Split(query, "&")
	queryDict := NewQueryDict()
	for _, part := range fields {
		if strings.Contains(part, "=") {
			subs := strings.Split(part, "=")
			queryDict.Set(subs[0], strings.Split(subs[1], ","))
		}
	}
	return queryDict
}

func (qdict *QueryDict) Set(key string, values ...interface{}) {
	qdict.dict[key] = values
}

func (qdict *QueryDict) SetAttr(key string, value interface{}) {
	if vals, exists := qdict.dict[key]; exists {
		qdict.dict[key] = append(vals, value)
	} else {
		qdict.dict[key] = []interface{}{value}
	}
}

func (qdict QueryDict) Get(key string, defaults ...interface{}) []interface{} {
	if values, exists := qdict.dict[key]; exists {
		return values
	} else if len(defaults) > 0 {
		return defaults
	}
	return nil
}

func (qdict QueryDict) Clone() *QueryDict {
	new_qdict := qdict
	return &new_qdict
}

type RouteStack struct {
	_data []*Route
}

type RouteStackInterface interface {
	IsEmpty() bool
	Size() int
	Top() *Route
	Push(data *Route)
	PushArray(arr []*Route)
	Pop() *Route
	get(index int) *Route
}

func NewRouteStack() *RouteStack {
	return &RouteStack{_data: make([]*Route, 0)}
}

func (RouteStack RouteStack) Size() int {
	return len(RouteStack._data)
}

func (RouteStack RouteStack) IsEmpty() bool {
	return RouteStack.Size() == 0
}

func (RouteStack RouteStack) get(index int) (*Route, error) {
	if index < RouteStack.Size() {
		return RouteStack._data[index], nil
	}
	return nil, errors.New("index " + strconv.Itoa(index) + " out of range")
}

func (RouteStack RouteStack) Top() *Route {
	if data, err := RouteStack.get(RouteStack.Size() - 1); err == nil {
		return data
	}
	return nil
}

func (RouteStack *RouteStack) Push(data *Route) {
	RouteStack._data = append(RouteStack._data, data)
}
func (RouteStack *RouteStack) PushArray(arr []*Route) {
	for i := len(arr) - 1; i >= 0; i = i - 1 {
		RouteStack.Push(arr[i])
	}
}

func (RouteStack *RouteStack) Pop() *Route {
	data := RouteStack.Top()
	if data != nil {
		RouteStack._data = RouteStack._data[0 : RouteStack.Size()-1]
	}
	return data
}
