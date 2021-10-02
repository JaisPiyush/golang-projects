package pkg

import "strings"

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
