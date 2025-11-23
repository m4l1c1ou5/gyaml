package gyaml

import (
	"testing"
)

const testYAML = `
name:
  first: Tom
  last: Anderson
age: 37
children:
  - Sara
  - Alex
  - Jack
fav.movie: Deer Hunter
friends:
  - first: Dale
    last: Murphy
    age: 44
    nets:
      - ig
      - fb
      - tw
  - first: Roger
    last: Craig
    age: 68
    nets:
      - fb
      - tw
  - first: Jane
    last: Murphy
    age: 47
    nets:
      - ig
      - tw
`

func TestGet(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"name.last", "Anderson"},
		{"age", "37"},
		{"children.0", "Sara"},
		{"children.1", "Alex"},
		{"children.2", "Jack"},
		{"friends.0.first", "Dale"},
		{"friends.1.last", "Craig"},
	}

	for _, tt := range tests {
		result := Get(testYAML, tt.path)
		if result.String() != tt.expected {
			t.Errorf("Get(%q) = %q, want %q", tt.path, result.String(), tt.expected)
		}
	}
}

func TestGetCount(t *testing.T) {
	result := Get(testYAML, "children.#")
	if result.Int() != 3 {
		t.Errorf("children.# = %d, want 3", result.Int())
	}

	result = Get(testYAML, "friends.#")
	if result.Int() != 3 {
		t.Errorf("friends.# = %d, want 3", result.Int())
	}
}

func TestGetArray(t *testing.T) {
	result := Get(testYAML, "children")
	arr := result.Array()
	if len(arr) != 3 {
		t.Errorf("len(children) = %d, want 3", len(arr))
	}

	expected := []string{"Sara", "Alex", "Jack"}
	for i, item := range arr {
		if item.String() != expected[i] {
			t.Errorf("children[%d] = %q, want %q", i, item.String(), expected[i])
		}
	}
}

func TestGetNestedArray(t *testing.T) {
	result := Get(testYAML, "friends.#.first")
	arr := result.Array()

	expected := []string{"Dale", "Roger", "Jane"}
	if len(arr) != len(expected) {
		t.Errorf("len(friends.#.first) = %d, want %d", len(arr), len(expected))
		return
	}

	for i, item := range arr {
		if item.String() != expected[i] {
			t.Errorf("friends.#.first[%d] = %q, want %q", i, item.String(), expected[i])
		}
	}
}

func TestQuery(t *testing.T) {
	result := Get(testYAML, `friends.#(last=="Murphy").first`)
	if result.String() != "Dale" {
		t.Errorf("query result = %q, want %q", result.String(), "Dale")
	}
}

func TestQueryMulti(t *testing.T) {
	result := Get(testYAML, `friends.#(last=="Murphy")#.first`)
	arr := result.Array()

	expected := []string{"Dale", "Jane"}
	if len(arr) != len(expected) {
		t.Errorf("len(query result) = %d, want %d", len(arr), len(expected))
		return
	}

	for i, item := range arr {
		if item.String() != expected[i] {
			t.Errorf("query result[%d] = %q, want %q", i, item.String(), expected[i])
		}
	}
}

func TestQueryComparison(t *testing.T) {
	result := Get(testYAML, "friends.#(age>45)#.last")
	arr := result.Array()

	// Should match Roger (68) and Jane (47)
	if len(arr) < 2 {
		t.Errorf("len(age>45) = %d, want at least 2", len(arr))
	}
}

func TestTypes(t *testing.T) {
	// String
	result := Get(testYAML, "name.first")
	if result.Type != String {
		t.Errorf("name.first type = %v, want String", result.Type)
	}
	if result.String() != "Tom" {
		t.Errorf("name.first = %q, want %q", result.String(), "Tom")
	}

	// Number
	result = Get(testYAML, "age")
	if result.Type != Number {
		t.Errorf("age type = %v, want Number", result.Type)
	}
	if result.Int() != 37 {
		t.Errorf("age = %d, want 37", result.Int())
	}
}

func TestExists(t *testing.T) {
	if !Get(testYAML, "name.last").Exists() {
		t.Error("name.last should exist")
	}

	if Get(testYAML, "name.middle").Exists() {
		t.Error("name.middle should not exist")
	}

	if Get(testYAML, "invalid.path.here").Exists() {
		t.Error("invalid.path.here should not exist")
	}
}

func TestIsArray(t *testing.T) {
	if !Get(testYAML, "children").IsArray() {
		t.Error("children should be an array")
	}

	if Get(testYAML, "name").IsArray() {
		t.Error("name should not be an array")
	}
}

func TestIsObject(t *testing.T) {
	if !Get(testYAML, "name").IsObject() {
		t.Error("name should be an object")
	}

	if Get(testYAML, "children").IsObject() {
		t.Error("children should not be an object")
	}
}

func TestMap(t *testing.T) {
	result := Get(testYAML, "name")
	m := result.Map()

	if m["first"].String() != "Tom" {
		t.Errorf("name.first = %q, want %q", m["first"].String(), "Tom")
	}

	if m["last"].String() != "Anderson" {
		t.Errorf("name.last = %q, want %q", m["last"].String(), "Anderson")
	}
}

func TestForEach(t *testing.T) {
	result := Get(testYAML, "children")

	var items []string
	result.ForEach(func(key, value Result) bool {
		items = append(items, value.String())
		return true
	})

	expected := []string{"Sara", "Alex", "Jack"}
	if len(items) != len(expected) {
		t.Errorf("len(items) = %d, want %d", len(items), len(expected))
		return
	}

	for i, item := range items {
		if item != expected[i] {
			t.Errorf("items[%d] = %q, want %q", i, item, expected[i])
		}
	}
}

func TestForEachObject(t *testing.T) {
	result := Get(testYAML, "name")

	m := make(map[string]string)
	result.ForEach(func(key, value Result) bool {
		m[key.String()] = value.String()
		return true
	})

	if m["first"] != "Tom" {
		t.Errorf("first = %q, want %q", m["first"], "Tom")
	}

	if m["last"] != "Anderson" {
		t.Errorf("last = %q, want %q", m["last"], "Anderson")
	}
}

func TestParse(t *testing.T) {
	result := Parse(testYAML)
	if result.Type != YAML {
		t.Errorf("Parse type = %v, want YAML", result.Type)
	}

	// Should be able to get values from parsed result
	name := result.Get("name.first")
	if name.String() != "Tom" {
		t.Errorf("parsed name.first = %q, want %q", name.String(), "Tom")
	}
}

func TestValid(t *testing.T) {
	if !Valid(testYAML) {
		t.Error("testYAML should be valid")
	}

	if Valid("invalid: yaml: content: [") {
		t.Error("invalid yaml should not be valid")
	}

	if Valid("") {
		t.Error("empty string should not be valid")
	}
}

func TestGetMany(t *testing.T) {
	result := GetMany(testYAML, "name.first", "age", "children.0")
	arr := result.Array()

	if len(arr) != 3 {
		t.Errorf("len(GetMany) = %d, want 3", len(arr))
		return
	}

	expected := []string{"Tom", "37", "Sara"}
	for i, item := range arr {
		if item.String() != expected[i] {
			t.Errorf("GetMany[%d] = %q, want %q", i, item.String(), expected[i])
		}
	}
}

func TestModifierReverse(t *testing.T) {
	result := Get(testYAML, "children|@reverse")
	arr := result.Array()

	expected := []string{"Jack", "Alex", "Sara"}
	if len(arr) != len(expected) {
		t.Errorf("len(@reverse) = %d, want %d", len(arr), len(expected))
		return
	}

	for i, item := range arr {
		if item.String() != expected[i] {
			t.Errorf("@reverse[%d] = %q, want %q", i, item.String(), expected[i])
		}
	}
}

func TestModifierKeys(t *testing.T) {
	result := Get(testYAML, "name|@keys")
	arr := result.Array()

	if len(arr) != 2 {
		t.Errorf("len(@keys) = %d, want 2", len(arr))
		return
	}

	// Keys should be "first" and "last" (order may vary)
	keys := make(map[string]bool)
	for _, item := range arr {
		keys[item.String()] = true
	}

	if !keys["first"] || !keys["last"] {
		t.Error("@keys should contain 'first' and 'last'")
	}
}

func TestModifierValues(t *testing.T) {
	result := Get(testYAML, "name|@values")
	arr := result.Array()

	if len(arr) != 2 {
		t.Errorf("len(@values) = %d, want 2", len(arr))
		return
	}

	// Values should be "Tom" and "Anderson" (order may vary)
	values := make(map[string]bool)
	for _, item := range arr {
		values[item.String()] = true
	}

	if !values["Tom"] || !values["Anderson"] {
		t.Error("@values should contain 'Tom' and 'Anderson'")
	}
}

func TestModifierThis(t *testing.T) {
	result := Get(testYAML, "@this")
	if result.Type != YAML {
		t.Errorf("@this type = %v, want YAML", result.Type)
	}
}

func TestModifierValid(t *testing.T) {
	result := Get(testYAML, "@valid")
	if result.String() != "true" {
		t.Errorf("@valid = %q, want %q", result.String(), "true")
	}
}

func TestBoolConversion(t *testing.T) {
	yaml := `
enabled: true
disabled: false
`

	if !Get(yaml, "enabled").Bool() {
		t.Error("enabled should be true")
	}

	if Get(yaml, "disabled").Bool() {
		t.Error("disabled should be false")
	}
}

func TestIntConversion(t *testing.T) {
	yaml := `
count: 42
negative: -10
`

	if Get(yaml, "count").Int() != 42 {
		t.Errorf("count = %d, want 42", Get(yaml, "count").Int())
	}

	if Get(yaml, "negative").Int() != -10 {
		t.Errorf("negative = %d, want -10", Get(yaml, "negative").Int())
	}
}

func TestFloatConversion(t *testing.T) {
	yaml := `
pi: 3.14159
negative: -2.5
`

	pi := Get(yaml, "pi").Float()
	if pi < 3.14 || pi > 3.15 {
		t.Errorf("pi = %f, want ~3.14159", pi)
	}

	if Get(yaml, "negative").Float() != -2.5 {
		t.Errorf("negative = %f, want -2.5", Get(yaml, "negative").Float())
	}
}

func TestValue(t *testing.T) {
	// String value
	result := Get(testYAML, "name.first")
	if v, ok := result.Value().(string); !ok || v != "Tom" {
		t.Errorf("Value() = %v, want %q", result.Value(), "Tom")
	}

	// Number value
	result = Get(testYAML, "age")
	if v, ok := result.Value().(float64); !ok || v != 37 {
		t.Errorf("Value() = %v, want 37", result.Value())
	}
}

func TestGetBytes(t *testing.T) {
	yamlBytes := []byte(testYAML)
	result := GetBytes(yamlBytes, "name.first")

	if result.String() != "Tom" {
		t.Errorf("GetBytes = %q, want %q", result.String(), "Tom")
	}
}

func TestEscapedKeys(t *testing.T) {
	yaml := `
fav.movie: Deer Hunter
`

	// This should work with escaped dot
	result := Get(yaml, "fav\\.movie")
	if result.String() != "Deer Hunter" {
		t.Errorf("escaped key = %q, want %q", result.String(), "Deer Hunter")
	}
}

func TestWildcard(t *testing.T) {
	result := Get(testYAML, "child*.2")
	if result.String() != "Jack" {
		t.Errorf("wildcard = %q, want %q", result.String(), "Jack")
	}
}

func TestEmptyPath(t *testing.T) {
	result := Get(testYAML, "")
	if result.Type != YAML {
		t.Errorf("empty path type = %v, want YAML", result.Type)
	}
}

func TestNonExistentPath(t *testing.T) {
	result := Get(testYAML, "does.not.exist")
	if result.Exists() {
		t.Error("non-existent path should not exist")
	}
	if result.Type != Null {
		t.Errorf("non-existent path type = %v, want Null", result.Type)
	}
}

func TestLess(t *testing.T) {
	a := Result{Type: Number, Num: 10}
	b := Result{Type: Number, Num: 20}

	if !a.Less(b, false) {
		t.Error("10 should be less than 20")
	}

	if b.Less(a, false) {
		t.Error("20 should not be less than 10")
	}
}

func TestAddModifier(t *testing.T) {
	AddModifier("test", func(yaml, arg string) string {
		return "modified"
	})

	result := Get(testYAML, "@test")
	if result.String() != "modified" {
		t.Errorf("custom modifier = %q, want %q", result.String(), "modified")
	}
}

// Benchmark tests

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(testYAML, "name.last")
	}
}

func BenchmarkGetNested(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(testYAML, "friends.0.first")
	}
}

func BenchmarkGetArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(testYAML, "children.#.0")
	}
}

func BenchmarkQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(testYAML, `friends.#(last=="Murphy").first`)
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(testYAML)
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid(testYAML)
	}
}
