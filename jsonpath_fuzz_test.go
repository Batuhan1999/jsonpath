package jsonpath

import (
	"encoding/json"
	"testing"
)

// A '*' token appearing when no tokens have been collected yet used to index
// tokens[-1] and panic. Parsing must not panic on any input.
func TestTokenize_StarNoPanic(t *testing.T) {
	for _, p := range []string{`""*`, `*`, `$*`, `**`, `$..*`, `$.a.*`} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Compile(%q) panicked: %v", p, r)
				}
			}()
			_, _ = Compile(p)
		}()
	}
}

func FuzzJsonPathLookup(f *testing.F) {
	var data interface{}
	_ = json.Unmarshal([]byte(`{"a":[{"b":1},{"b":2}],"c":{"d":3},"e":[1,2,3]}`), &data)
	for _, p := range []string{"$.a", "$.a[0].b", "$..b", "$.e[0:2]", "$.a[*]", "$.e[-1]", "$.a[?(@.b>1)]", "$", `""*`, "$..*"} {
		f.Add(p)
	}
	f.Fuzz(func(_ *testing.T, path string) {
		_, _ = JsonPathLookup(data, path)
	})
}

// A reversed or out-of-range slice (e.g. [1:0]) must return an empty result,
// not panic in reflect.Value.Slice.
func TestGetRange_NoPanic(t *testing.T) {
	var data interface{}
	_ = json.Unmarshal([]byte(`{"e":[1,2,3]}`), &data)
	for _, p := range []string{"$.e[1:0]", "$.e[2:1]", "$.e[1:100]", "$.e[-1:0]", "$e[1:00"} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("JsonPathLookup(%q) panicked: %v", p, r)
				}
			}()
			_, _ = JsonPathLookup(data, p)
		}()
	}
}
