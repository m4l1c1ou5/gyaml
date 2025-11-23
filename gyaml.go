// Copyright 2024 GYAML Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package gyaml provides searching for yaml strings.
package gyaml

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	yamlv3 "gopkg.in/yaml.v3"
)

// Type is Result type
type Type int

const (
	// Null is a null yaml value
	Null Type = iota
	// False is a yaml false boolean
	False
	// Number is yaml number
	Number
	// String is a yaml string
	String
	// True is a yaml true boolean
	True
	// YAML is a raw block of YAML
	YAML
)

// String returns a string representation of the type.
func (t Type) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case YAML:
		return "YAML"
	}
}

// Result represents a yaml value that is returned from Get().
type Result struct {
	// Type is the yaml type
	Type Type
	// Raw is the raw yaml
	Raw string
	// Str is the yaml string
	Str string
	// Num is the yaml number
	Num float64
	// Index of raw value in original yaml, zero means index unknown
	Index int
	// Indexes of all the elements that match on a path containing the '#'
	// query character.
	Indexes []int
}

// String returns a string representation of the value.
func (t Result) String() string {
	switch t.Type {
	case Number:
		if len(t.Raw) == 0 {
			// calculated result
			return strconv.FormatFloat(t.Num, 'f', -1, 64)
		}
		var i int
		if t.Raw[0] == '-' {
			i++
		}
		for ; i < len(t.Raw); i++ {
			if t.Raw[i] < '0' || t.Raw[i] > '9' {
				return strconv.FormatFloat(t.Num, 'f', -1, 64)
			}
		}
		return t.Raw
	case String:
		return t.Str
	case YAML:
		return t.Raw
	case True:
		return "true"
	case False:
		return "false"
	default:
		return ""
	}
}

// Bool returns a boolean representation.
func (t Result) Bool() bool {
	switch t.Type {
	default:
		return false
	case True:
		return true
	case String:
		b, _ := strconv.ParseBool(strings.ToLower(t.Str))
		return b
	case Number:
		return t.Num != 0
	}
}

// Int returns an integer representation.
func (t Result) Int() int64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseInt(t.Str)
		return n
	case Number:
		// try to directly convert the float64 to int64
		i, ok := safeInt(t.Num)
		if ok {
			return i
		}
		// now try to parse the raw string
		i, ok = parseInt(t.Raw)
		if ok {
			return i
		}
		// fallback to a standard conversion
		return int64(t.Num)
	}
}

// Uint returns an unsigned integer representation.
func (t Result) Uint() uint64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := parseUint(t.Str)
		return n
	case Number:
		// try to directly convert the float64 to uint64
		i, ok := safeInt(t.Num)
		if ok && i >= 0 {
			return uint64(i)
		}
		// now try to parse the raw string
		u, ok := parseUint(t.Raw)
		if ok {
			return u
		}
		// fallback to a standard conversion
		return uint64(t.Num)
	}
}

// Float returns a float64 representation.
func (t Result) Float() float64 {
	switch t.Type {
	default:
		return 0
	case True:
		return 1
	case String:
		n, _ := strconv.ParseFloat(t.Str, 64)
		return n
	case Number:
		return t.Num
	}
}

// Time returns a time.Time representation.
func (t Result) Time() time.Time {
	res, _ := time.Parse(time.RFC3339, t.String())
	return res
}

// Array returns back an array of values.
// If the result represents a null value or is non-existent, then an empty
// array will be returned.
// If the result is not a YAML array, the return value will be an
// array containing one result.
func (t Result) Array() []Result {
	if t.Type == Null {
		return []Result{}
	}
	if !t.IsArray() {
		return []Result{t}
	}
	r := t.arrayOrMap('[', false)
	return r.a
}

// IsObject returns true if the result value is a YAML object.
func (t Result) IsObject() bool {
	return t.Type == YAML && len(t.Raw) > 0 && (t.Raw[0] == '{' || isYAMLObject(t.Raw))
}

// IsArray returns true if the result value is a YAML array.
func (t Result) IsArray() bool {
	return t.Type == YAML && len(t.Raw) > 0 && (t.Raw[0] == '[' || isYAMLArray(t.Raw))
}

// IsBool returns true if the result value is a YAML boolean.
func (t Result) IsBool() bool {
	return t.Type == True || t.Type == False
}

// ForEach iterates through values.
// If the result represents a non-existent value, then no values will be
// iterated. If the result is an Object, the iterator will pass the key and
// value of each item. If the result is an Array, the iterator will only pass
// the value of each item. If the result is not a YAML array or object, the
// iterator will pass back one value equal to the result.
func (t Result) ForEach(iterator func(key, value Result) bool) {
	if !t.Exists() {
		return
	}
	if t.Type != YAML {
		iterator(Result{}, t)
		return
	}

	// Parse the YAML to determine structure
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(t.Raw), &data); err != nil {
		return
	}

	switch v := data.(type) {
	case map[string]interface{}:
		// Object iteration
		for key, val := range v {
			keyResult := Result{Type: String, Str: key, Raw: key}
			valResult := valueToResult(val)
			if !iterator(keyResult, valResult) {
				return
			}
		}
	case []interface{}:
		// Array iteration
		for i, val := range v {
			keyResult := Result{Type: Number, Num: float64(i)}
			valResult := valueToResult(val)
			if !iterator(keyResult, valResult) {
				return
			}
		}
	default:
		// Single value
		iterator(Result{}, t)
	}
}

// Map returns back a map of values. The result should be a YAML object.
// If the result is not a YAML object, the return value will be an empty map.
func (t Result) Map() map[string]Result {
	if t.Type != YAML {
		return map[string]Result{}
	}
	r := t.arrayOrMap('{', false)
	return r.o
}

// Get searches result for the specified path.
// The result should be a YAML array or object.
func (t Result) Get(path string) Result {
	r := Get(t.Raw, path)
	if r.Indexes != nil {
		for i := 0; i < len(r.Indexes); i++ {
			r.Indexes[i] += t.Index
		}
	} else {
		r.Index += t.Index
	}
	return r
}

type arrayOrMapResult struct {
	a  []Result
	ai []interface{}
	o  map[string]Result
	oi map[string]interface{}
	vc byte
}

func (t Result) arrayOrMap(vc byte, valueize bool) (r arrayOrMapResult) {
	// Parse YAML to get structure
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(t.Raw), &data); err != nil {
		return
	}

	switch v := data.(type) {
	case map[string]interface{}:
		r.vc = '{'
		if valueize {
			r.oi = v
		} else {
			r.o = make(map[string]Result)
			for key, val := range v {
				r.o[key] = valueToResult(val)
			}
		}
	case []interface{}:
		r.vc = '['
		if valueize {
			r.ai = v
		} else {
			r.a = make([]Result, len(v))
			for i, val := range v {
				r.a[i] = valueToResult(val)
			}
		}
	}

	return r
}

// Exists returns true if value exists.
//
//	if gjson.Get(yaml, "name.last").Exists() {
//		println("has a last name")
//	}
func (t Result) Exists() bool {
	return t.Type != Null || len(t.Raw) != 0
}

// Value returns one of these types:
//
//	bool, for YAML booleans
//	float64, for YAML numbers
//	string, for YAML strings
//	nil, for YAML null
//	[]interface{}, for YAML arrays
//	map[string]interface{}, for YAML objects
func (t Result) Value() interface{} {
	switch t.Type {
	default:
		return nil
	case False:
		return false
	case Number:
		return t.Num
	case String:
		return t.Str
	case True:
		return true
	case YAML:
		var data interface{}
		if err := yamlv3.Unmarshal([]byte(t.Raw), &data); err != nil {
			return nil
		}
		return data
	}
}

// Less return true if a token is less than another token.
// The caseSensitive parameter is used when the tokens are Strings.
// The order when comparing two different type is:
//
//	Null < False < Number < String < True < YAML
func (t Result) Less(token Result, caseSensitive bool) bool {
	if t.Type < token.Type {
		return true
	}
	if t.Type > token.Type {
		return false
	}
	if t.Type == String {
		if caseSensitive {
			return t.Str < token.Str
		}
		return stringLessInsensitive(t.Str, token.Str)
	}
	if t.Type == Number {
		return t.Num < token.Num
	}
	return t.Raw < token.Raw
}

func stringLessInsensitive(a, b string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] >= 'A' && a[i] <= 'Z' {
			if b[i] >= 'A' && b[i] <= 'Z' {
				if a[i] < b[i] {
					return true
				} else if a[i] > b[i] {
					return false
				}
			} else {
				if a[i]+32 < b[i] {
					return true
				} else if a[i]+32 > b[i] {
					return false
				}
			}
		} else if b[i] >= 'A' && b[i] <= 'Z' {
			if a[i] < b[i]+32 {
				return true
			} else if a[i] > b[i]+32 {
				return false
			}
		} else {
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
	}
	return len(a) < len(b)
}

// parseContext is used during parsing
type parseContext struct {
	yaml   string
	value  Result
	pipe   string
	piped  bool
	calcd  bool
	lines  bool
	yamlma map[string]interface{}
}

// Get searches yaml for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// This function expects that the yaml is well-formed and validates. Bad yaml will not panic,
// but it may return unexpected results.
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "fav.movie": "Deer Hunter"
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children.0"         >> "Sara"
//	"children.1"         >> "Alex"
func Get(yaml, path string) Result {
	if len(path) > 1 && path[0] == '.' && path[1] == '.' {
		return getMany(yaml, path)
	}

	if len(path) == 0 {
		// empty path returns the entire yaml
		return Result{
			Type:  YAML,
			Raw:   yaml,
			Index: 0,
		}
	}

	if path[0] == '.' {
		// path starts with dot, remove it
		path = path[1:]
	}

	// Try fast path first for simple queries
	if result, ok := fastGet(yaml, path); ok {
		return result
	}

	// Fall back to slow path for complex queries
	var c parseContext
	c.yaml = yaml

	// Convert YAML to a normalized form for easier parsing
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yaml), &data); err != nil {
		return c.value
	}
	c.yamlma = make(map[string]interface{})

	// Now traverse the path
	return getFromPath(data, path, yaml)
}

// GetBytes searches yaml for the specified path.
// If working with bytes, this method preferred over Get(string(data), path)
func GetBytes(yaml []byte, path string) Result {
	return Get(*(*string)(unsafe.Pointer(&yaml)), path)
}

// GetMany searches yaml for multiple paths.
// The return value is a Result holding a YAML array of values.
// An empty array is returned if the yaml is not valid.
func GetMany(yaml string, path ...string) Result {
	var res Result
	res.Type = YAML
	var data []byte
	data = append(data, '[')
	for i, path := range path {
		if i > 0 {
			data = append(data, ',')
		}
		val := Get(yaml, path)
		data = append(data, val.Raw...)
	}
	data = append(data, ']')
	res.Raw = string(data)
	return res
}

// GetManyBytes searches yaml for multiple paths.
// The return value is a Result holding a YAML array of values.
func GetManyBytes(yaml []byte, path ...string) Result {
	return GetMany(string(yaml), path...)
}

func getMany(yaml, path string) Result {
	// Handle lines (..) prefix
	var data []interface{}
	lines := strings.Split(yaml, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var item interface{}
		if err := yamlv3.Unmarshal([]byte(line), &item); err != nil {
			continue
		}
		data = append(data, item)
	}

	// Remove the .. prefix
	path = path[2:]
	return getFromPath(data, path, yaml)
}

// Parse parses the yaml and returns a result.
//
//	value := gyaml.Parse(yaml)
func Parse(yaml string) Result {
	var res Result
	res.Type = YAML
	res.Raw = yaml
	return res
}

// ParseBytes parses the yaml and returns a result.
// If working with bytes, this method preferred over Parse(string(data)).
func ParseBytes(yaml []byte) Result {
	return Parse(string(yaml))
}

// Valid returns true if the input is valid yaml.
//
//	if !gyaml.Valid(yaml) {
//		return errors.New("invalid yaml")
//	}
func Valid(yaml string) bool {
	if len(strings.TrimSpace(yaml)) == 0 {
		return false
	}
	var data interface{}
	return yamlv3.Unmarshal([]byte(yaml), &data) == nil
}

// ValidBytes returns true if the input is valid yaml.
func ValidBytes(yaml []byte) bool {
	return Valid(string(yaml))
}

// Helper functions

func parseInt(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err == nil
}

func parseUint(s string) (uint64, bool) {
	u, err := strconv.ParseUint(s, 10, 64)
	return u, err == nil
}

func safeInt(f float64) (int64, bool) {
	if f < -9007199254740991 || f > 9007199254740991 {
		return 0, false
	}
	i := int64(f)
	if float64(i) == f {
		return i, true
	}
	return 0, false
}

func valueToResult(val interface{}) Result {
	var res Result
	switch v := val.(type) {
	case nil:
		res.Type = Null
	case bool:
		if v {
			res.Type = True
			res.Raw = "true"
		} else {
			res.Type = False
			res.Raw = "false"
		}
	case int:
		res.Type = Number
		res.Num = float64(v)
		res.Raw = strconv.FormatInt(int64(v), 10)
	case int64:
		res.Type = Number
		res.Num = float64(v)
		res.Raw = strconv.FormatInt(v, 10)
	case float64:
		res.Type = Number
		res.Num = v
		res.Raw = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		res.Type = String
		res.Str = v
		res.Raw = v
	case []interface{}, map[string]interface{}:
		res.Type = YAML
		data, _ := yamlv3.Marshal(v)
		res.Raw = string(data)
	default:
		res.Type = YAML
		data, _ := yamlv3.Marshal(v)
		res.Raw = string(data)
	}
	return res
}

func isYAMLObject(s string) bool {
	// Simple heuristic: if it contains ": " it's likely an object
	return strings.Contains(s, ": ")
}

func isYAMLArray(s string) bool {
	// Simple heuristic: if it starts with "- " it's likely an array
	return strings.HasPrefix(strings.TrimSpace(s), "- ")
}

// getFromPath traverses a parsed YAML structure using a path
func getFromPath(data interface{}, path string, origYAML string) Result {
	if path == "" || path == "@this" {
		return valueToResult(data)
	}

	// Handle modifiers
	if path[0] == '@' {
		return applyModifier(data, path, origYAML)
	}

	// Parse path components
	parts := parsePath(path)
	if len(parts) == 0 {
		return valueToResult(data)
	}

	return traversePath(data, parts, origYAML)
}

// pathComponent represents a single component of a path
type pathComponent struct {
	key     string
	isWild  bool
	isQuery bool
	query   string
	isIndex bool
	index   int
	isCount bool
	pipe    string
	hasPipe bool
	multi   bool // for #()# queries
}

func parsePath(path string) []pathComponent {
	var parts []pathComponent
	var current strings.Builder
	var inQuery bool
	var queryDepth int
	var escaped bool

	for i := 0; i < len(path); i++ {
		ch := path[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '#' && !inQuery {
			if i+1 < len(path) && path[i+1] == '(' {
				// Start of query
				inQuery = true
				queryDepth = 1
				i++ // skip the '('

				// Save current key if any
				if current.Len() > 0 {
					parts = append(parts, pathComponent{key: current.String()})
					current.Reset()
				}

				// Start collecting query
				parts = append(parts, pathComponent{isQuery: true})
				continue
			} else {
				// It's a count operation
				if current.Len() > 0 {
					parts = append(parts, pathComponent{key: current.String()})
					current.Reset()
				}
				parts = append(parts, pathComponent{isCount: true})
				continue
			}
		}

		if inQuery {
			if ch == '(' {
				queryDepth++
			} else if ch == ')' {
				queryDepth--
				if queryDepth == 0 {
					// End of query
					parts[len(parts)-1].query = current.String()
					current.Reset()
					inQuery = false

					// Check for multi query (#()#)
					if i+1 < len(path) && path[i+1] == '#' {
						parts[len(parts)-1].multi = true
						i++
					}
					continue
				}
			}
			current.WriteByte(ch)
			continue
		}

		if ch == '.' {
			if current.Len() > 0 {
				parts = append(parts, parseComponent(current.String()))
				current.Reset()
			}
			continue
		}

		if ch == '|' {
			// Pipe for modifiers
			if current.Len() > 0 {
				parts = append(parts, parseComponent(current.String()))
				current.Reset()
			}
			// Rest is pipe
			parts = append(parts, pathComponent{pipe: path[i+1:], hasPipe: true})
			break
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		parts = append(parts, parseComponent(current.String()))
	}

	return parts
}

func parseComponent(s string) pathComponent {
	var comp pathComponent

	// Check for wildcard
	if strings.ContainsAny(s, "*?") {
		comp.key = s
		comp.isWild = true
		return comp
	}

	// Check for index
	if idx, err := strconv.Atoi(s); err == nil {
		comp.isIndex = true
		comp.index = idx
		return comp
	}

	comp.key = s
	return comp
}

func traversePath(data interface{}, parts []pathComponent, origYAML string) Result {
	current := data

	for i, part := range parts {
		if part.hasPipe {
			// Apply pipe modifier
			res := valueToResult(current)
			return applyModifier(current, part.pipe, res.Raw)
		}

		if part.isCount {
			// Count operation - but check if there are more parts after this
			if i+1 < len(parts) {
				// There are more parts, so # means "apply to all elements"
				switch v := current.(type) {
				case []interface{}:
					// Apply remaining path to all elements
					var results []interface{}
					for _, item := range v {
						res := traversePath(item, parts[i+1:], origYAML)
						if res.Exists() {
							// Extract the actual value
							results = append(results, res.Value())
						}
					}
					return valueToResult(results)
				case map[string]interface{}:
					// Can't iterate over map with #
					return Result{Type: Null}
				}
			} else {
				// Just return count
				switch v := current.(type) {
				case []interface{}:
					return Result{Type: Number, Num: float64(len(v)), Raw: strconv.Itoa(len(v))}
				case map[string]interface{}:
					return Result{Type: Number, Num: float64(len(v)), Raw: strconv.Itoa(len(v))}
				default:
					return Result{Type: Number, Num: 0, Raw: "0"}
				}
			}
		}

		if part.isQuery {
			// Handle query
			current = handleQuery(current, part, parts[i+1:])
			if !part.multi {
				// Single match, continue with remaining path
				if i+1 < len(parts) {
					return traversePath(current, parts[i+1:], origYAML)
				}
			} else {
				// Multi match - if there are remaining parts, apply them to each match
				if i+1 < len(parts) {
					if matches, ok := current.([]interface{}); ok {
						var results []interface{}
						for _, match := range matches {
							res := traversePath(match, parts[i+1:], origYAML)
							if res.Exists() {
								results = append(results, res.Value())
							}
						}
						return valueToResult(results)
					}
				}
			}
			return valueToResult(current)
		}

		switch v := current.(type) {
		case map[string]interface{}:
			if part.isWild {
				// Wildcard match on object keys
				var matches []interface{}
				for key, val := range v {
					if matchPattern(key, part.key) {
						matches = append(matches, val)
					}
				}
				if len(matches) == 1 {
					current = matches[0]
				} else {
					current = matches
				}
			} else {
				val, ok := v[part.key]
				if !ok {
					return Result{Type: Null}
				}
				current = val
			}

		case []interface{}:
			if part.isIndex {
				if part.index < 0 || part.index >= len(v) {
					return Result{Type: Null}
				}
				current = v[part.index]
			} else if part.key != "" {
				// Apply to all elements in array
				var results []interface{}
				for _, item := range v {
					if m, ok := item.(map[string]interface{}); ok {
						if val, exists := m[part.key]; exists {
							results = append(results, val)
						}
					}
				}
				if len(results) == 0 {
					return Result{Type: Null}
				}
				current = results
			}

		default:
			return Result{Type: Null}
		}
	}

	return valueToResult(current)
}

func handleQuery(data interface{}, part pathComponent, remainingParts []pathComponent) interface{} {
	arr, ok := data.([]interface{})
	if !ok {
		return nil
	}

	query := part.query
	matches := []interface{}{}

	for _, item := range arr {
		if evaluateQuery(item, query) {
			matches = append(matches, item)
		}
	}

	if part.multi {
		// Return all matches
		return matches
	}

	// Return first match
	if len(matches) > 0 {
		return matches[0]
	}
	return nil
}

func evaluateQuery(item interface{}, query string) bool {
	// Parse query: key op value
	// Supported operators: ==, !=, <, <=, >, >=, %, !%

	var key, op, value string

	// Find operator
	for _, operator := range []string{"==", "!=", "<=", ">=", "<", ">", "!%", "%"} {
		if idx := strings.Index(query, operator); idx != -1 {
			key = strings.TrimSpace(query[:idx])
			op = operator
			value = strings.TrimSpace(query[idx+len(operator):])
			break
		}
	}

	if op == "" {
		return false
	}

	// Remove quotes from value
	value = strings.Trim(value, "\"'")

	// Get the value from item
	var itemValue interface{}
	if key == "" {
		itemValue = item
	} else {
		itemValue = getValueByKey(item, key)
	}

	return compareValues(itemValue, op, value)
}

func getValueByKey(data interface{}, key string) interface{} {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	// Handle nested keys
	if strings.Contains(key, ".") {
		parts := strings.Split(key, ".")
		current := data
		for _, part := range parts {
			if m, ok := current.(map[string]interface{}); ok {
				current = m[part]
			} else {
				return nil
			}
		}
		return current
	}

	return m[key]
}

func compareValues(itemValue interface{}, op string, value string) bool {
	switch op {
	case "==":
		return fmt.Sprint(itemValue) == value
	case "!=":
		return fmt.Sprint(itemValue) != value
	case "%":
		// Pattern match
		return matchPattern(fmt.Sprint(itemValue), value)
	case "!%":
		return !matchPattern(fmt.Sprint(itemValue), value)
	case "<", "<=", ">", ">=":
		return compareNumeric(itemValue, op, value)
	}
	return false
}

func compareNumeric(itemValue interface{}, op string, value string) bool {
	var itemNum float64
	switch v := itemValue.(type) {
	case float64:
		itemNum = v
	case int:
		itemNum = float64(v)
	case int64:
		itemNum = float64(v)
	case string:
		itemNum, _ = strconv.ParseFloat(v, 64)
	default:
		return false
	}

	valueNum, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}

	switch op {
	case "<":
		return itemNum < valueNum
	case "<=":
		return itemNum <= valueNum
	case ">":
		return itemNum > valueNum
	case ">=":
		return itemNum >= valueNum
	}
	return false
}

func matchPattern(str, pattern string) bool {
	return wildcard(str, pattern)
}

// wildcard implements wildcard pattern matching
func wildcard(str, pattern string) bool {
	return deepMatch(str, pattern)
}

func deepMatch(str, pattern string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			if len(pattern) == 1 {
				return true
			}
			for i := 0; i <= len(str); i++ {
				if deepMatch(str[i:], pattern[1:]) {
					return true
				}
			}
			return false
		case '?':
			if len(str) == 0 {
				return false
			}
			str = str[1:]
			pattern = pattern[1:]
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
			str = str[1:]
			pattern = pattern[1:]
		}
	}
	return len(str) == 0
}

// Modifiers

var modifiers = map[string]func(yaml, arg string) string{
	"reverse": modReverse,
	"ugly":    modUgly,
	"pretty":  modPretty,
	"this":    modThis,
	"valid":   modValid,
	"flatten": modFlatten,
	"join":    modJoin,
	"keys":    modKeys,
	"values":  modValues,
}

// AddModifier adds a custom modifier
func AddModifier(name string, fn func(yaml, arg string) string) {
	modifiers[name] = fn
}

func applyModifier(data interface{}, path string, yamlStr string) Result {
	if !strings.HasPrefix(path, "@") {
		return valueToResult(data)
	}

	// Parse modifier name and argument
	var modName, modArg string
	if idx := strings.Index(path, ":"); idx != -1 {
		modName = path[1:idx]
		modArg = path[idx+1:]
	} else {
		modName = path[1:]
	}

	// Check for pipe after modifier
	if idx := strings.Index(modName, "|"); idx != -1 {
		remaining := modName[idx:]
		modName = modName[:idx]
		// Apply modifier then continue with pipe
		if fn, ok := modifiers[modName]; ok {
			result := fn(yamlStr, modArg)
			return applyModifier(data, remaining[1:], result)
		}
	}

	if fn, ok := modifiers[modName]; ok {
		result := fn(yamlStr, modArg)
		var newData interface{}
		yamlv3.Unmarshal([]byte(result), &newData)
		return valueToResult(newData)
	}

	return valueToResult(data)
}

func modReverse(yamlStr, arg string) string {
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}

	switch v := data.(type) {
	case []interface{}:
		reversed := make([]interface{}, len(v))
		for i, item := range v {
			reversed[len(v)-1-i] = item
		}
		result, _ := yamlv3.Marshal(reversed)
		return string(result)
	case map[string]interface{}:
		// Reverse doesn't really apply to maps in the same way
		result, _ := yamlv3.Marshal(v)
		return string(result)
	}
	return yamlStr
}

func modUgly(yamlStr, arg string) string {
	// Remove unnecessary whitespace
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}
	result, _ := yamlv3.Marshal(data)
	return strings.TrimSpace(string(result))
}

func modPretty(yamlStr, arg string) string {
	// YAML is already pretty by default
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}
	result, _ := yamlv3.Marshal(data)
	return string(result)
}

func modThis(yamlStr, arg string) string {
	return yamlStr
}

func modValid(yamlStr, arg string) string {
	if Valid(yamlStr) {
		return "true"
	}
	return "false"
}

func modFlatten(yamlStr, arg string) string {
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}

	if arr, ok := data.([]interface{}); ok {
		flattened := flattenArray(arr)
		result, _ := yamlv3.Marshal(flattened)
		return string(result)
	}
	return yamlStr
}

func flattenArray(arr []interface{}) []interface{} {
	var result []interface{}
	for _, item := range arr {
		if subArr, ok := item.([]interface{}); ok {
			result = append(result, flattenArray(subArr)...)
		} else {
			result = append(result, item)
		}
	}
	return result
}

func modJoin(yamlStr, arg string) string {
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}

	if arr, ok := data.([]interface{}); ok {
		joined := make(map[string]interface{})
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				for k, v := range m {
					joined[k] = v
				}
			}
		}
		result, _ := yamlv3.Marshal(joined)
		return string(result)
	}
	return yamlStr
}

func modKeys(yamlStr, arg string) string {
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}

	if m, ok := data.(map[string]interface{}); ok {
		keys := make([]interface{}, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		result, _ := yamlv3.Marshal(keys)
		return string(result)
	}
	return yamlStr
}

func modValues(yamlStr, arg string) string {
	var data interface{}
	if err := yamlv3.Unmarshal([]byte(yamlStr), &data); err != nil {
		return yamlStr
	}

	if m, ok := data.(map[string]interface{}); ok {
		values := make([]interface{}, 0, len(m))
		for _, v := range m {
			values = append(values, v)
		}
		result, _ := yamlv3.Marshal(values)
		return string(result)
	}
	return yamlStr
}

// ForEachLine iterates through lines of YAML
func ForEachLine(yaml string, iterator func(line Result) bool) {
	lines := strings.Split(yaml, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		res := Parse(line)
		if !iterator(res) {
			return
		}
	}
}
