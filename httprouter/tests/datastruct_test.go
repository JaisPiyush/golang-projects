package tests

import (
	"fmt"
	"reflect"
	"testing"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func TestQueryDict(t *testing.T) {
	queryDict := pkg.NewQueryDict()

	queryDict.SetAttr("key1", "value1")
	if value := queryDict.Get("key1"); value == nil {
		t.Error("QueryDict.SetAttr is not working")
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		t.Error("QueryDict.SetAttr is not intializing array")
	} else {
		fmt.Println(queryDict)
	}

	queryDict.Set("arome", "ar", 1, false)
	if val := queryDict.Get("arome"); len(val) != 3 {
		t.Error("vardiac params is not working in QueryDict.Set")
	}

	queryDict.Set("aroma", []interface{}{"a", "b", "d"}...)
	if val := queryDict.Get("aroma"); len(val) == 1 {
		t.Error("Array Destruction not working in QueryDict.Set")
	}

	val := queryDict.Get("mike", "test")
	if val == nil || val[0] != "test" {
		t.Error("defaults is not working in QueryDict.Get")
	}

	clone := queryDict.Clone()
	if clone == queryDict && &clone == &queryDict {
		t.Error("QueryDict.Clone not working")
	}

}
