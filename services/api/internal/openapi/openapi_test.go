package openapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
)

func TestGridworksOpenAPISemantics(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "packages", "schemas", "openapi", "gridworks-api-v1.json")
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	actual := map[string][]string{}
	for path, item := range raw.Paths {
		for method := range item {
			actual[path] = append(actual[path], strings.ToUpper(method))
		}
	}
	for path := range actual {
		sort.Strings(actual[path])
	}
	expected := api.RouteInventory()
	for path := range expected {
		sort.Strings(expected[path])
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("OpenAPI/handler inventory drift: %v != %v", actual, expected)
	}
}
