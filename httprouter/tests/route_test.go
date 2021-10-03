package tests

import (
	"testing"

	"gihtub.com/JaisPiyush/golang-projects/pkg"
)

func TestParsePart(t *testing.T) {

	type expectedOutput struct {
		name       string
		is_dynamic bool
		regex      string
	}

	patterns := []struct {
		input  string
		output expectedOutput
	}{
		{
			input: "",
			output: expectedOutput{
				name:       "",
				is_dynamic: false,
				regex:      "",
			},
		},
		{
			input: "name",
			output: expectedOutput{
				name:       "name",
				is_dynamic: false,
				regex:      "",
			},
		},
		{
			input: "{key}",
			output: expectedOutput{
				name:       "key",
				is_dynamic: true,
				regex:      "",
			},
		},
		{
			input: "{id:[0-9]+}",
			output: expectedOutput{
				name:       "id",
				is_dynamic: true,
				regex:      "[0-9]+",
			},
		},
	}

	for _, pattern := range patterns {
		route := pkg.ParsePart(pattern.input, "GET", nil)
		output := &pattern.output
		if !(route.Name == output.name && route.Regex == output.regex && route.IsDynamic == output.is_dynamic) {
			t.Error("Parsed Route is not producing result ", route.String(), output)
		}
	}
}

func TestIsRouteMatching(t *testing.T) {
	fixtures := []map[string]interface{}{
		{
			"part": "name",
			"route": &pkg.Route{Name: "name",
				Regex:     "",
				IsDynamic: false,
			},
			"pass": true,
		},
		{
			"part": "24",
			"route": &pkg.Route{
				Name:      "id",
				Regex:     "[0-9]+",
				IsDynamic: true,
			},
			"pass": true,
		},
		{
			"part": "piyush",
			"route": &pkg.Route{
				Name:      "id",
				Regex:     "[0-9]+",
				IsDynamic: true,
			},
			"pass": false,
		},
		{
			"part": "piyush",
			"route": &pkg.Route{
				Name:      "key",
				Regex:     "",
				IsDynamic: true,
			},
			"pass": true,
		},
	}

	for _, fixture := range fixtures {
		if matching := fixture["route"].(*pkg.Route).IsRouteMatching(fixture["part"].(string)); matching != fixture["pass"].(bool) {
			t.Error("Route Matching is failing", fixture["part"].(string), matching, fixture["pass"], fixture["route"].(*pkg.Route).String())
		}
	}
}
