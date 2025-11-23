# gyaml

Get YAML values quickly - YAML parser for Go

GYAML is a Go package that provides a fast and simple way to get values from a YAML document. It has `Get` fun, which takes YAML and a path, and returns a value.

Inspired by [gjson](https://github.com/tidwall/gjson) but for YAML documents.

## Getting Started

### Installing

To start using GYAML, install Go and run `go get`:

```sh
$ go get -u github.com/m4l1c1ou5/gyaml
```

This will retrieve the library.

### Get a value

Get searches YAML for the specified path. A path is in dot syntax, such as "name.last" or "age". When the value is found it's returned immediately.

```go
package main

import "github.com/m4l1c1ou5/gyaml"

const yaml = `
name:
  first: Janet
  last: Prichard
age: 47
`

func main() {
    value := gyaml.Get(yaml, "name.last")
    println(value.String())
}
```

This will print:

```
Prichard
```

There's also `GetBytes` for working with YAML byte slices.

### Path Syntax

Below is a quick overview of the path syntax, for more complete information please check out [GYAML Syntax](SYNTAX.md).

A path is a series of keys separated by a dot. A key may contain special wildcard characters '\*' and '?'. To access an array value use the index as the key. To get the number of elements in an array or to access a child path, use the '#' character. The dot and wildcard characters can be escaped with '\\'.

```yaml
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
```

```
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\\.movie"        >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```

You can also query an array for the first match by using `#(...)`, or find all matches with `#(...)#`. Queries support the `==`, `!=`, `<`, `<=`, `>`, `>=` comparison operators and the simple pattern matching `%` (like) and `!%` (not like) operators.

```
friends.#(last=="Murphy").first           >> "Dale"
friends.#(last=="Murphy")#.first          >> ["Dale","Jane"]
friends.#(age>45)#.last                   >> ["Craig","Murphy"]
friends.#(first%"D*").last                >> "Murphy"
friends.#(first!%"D*").last               >> "Craig"
friends.#(nets.#(=="fb"))#.first          >> ["Dale","Roger"]
```

## Result Type

GYAML supports the YAML types string, number, bool, and null. Arrays and Objects are returned as their raw YAML types.

The `Result` type holds one of these:

```
bool, for YAML booleans
float64, for YAML numbers
string, for YAML string literals
nil, for YAML null
```

To directly access the value:

```go
result.Type      // can be String, Number, True, False, Null, or YAML
result.Str       // holds the string
result.Num       // holds the float64 number
result.Raw       // holds the raw yaml
result.Index     // index of raw value in original yaml, zero means index unknown
result.Indexes   // indexes of all the elements that match on a path containing the '#' query character.
```

There are a variety of handy functions that work on a result:

```go
result.Exists() bool
result.Value() interface{}
result.Int() int64
result.Uint() uint64
result.Float() float64
result.String() string
result.Bool() bool
result.Time() time.Time
result.Array() []gyaml.Result
result.Map() map[string]gyaml.Result
result.Get(path string) Result
result.ForEach(iterator func(key, value Result) bool)
result.Less(token Result, caseSensitive bool) bool
```

The `result.Value()` function returns an `interface{}` which requires type assertion and is one of the following Go types:

```
boolean >> bool
number  >> float64
string  >> string
null    >> nil
array   >> []interface{}
object  >> map[string]interface{}
```

The `result.Array()` function returns back an array of values. If the result represents a non-existent value, then an empty array will be returned. If the result is not a YAML array, the return value will be an array containing one result.

### 64-bit integers

The `result.Int()` and `result.Uint()` calls are capable of reading all 64 bits, allowing for large integers that cannot be represented as a `float64`. For example, the number `9007199254740993` cannot be correctly represented as a `float64` but can be parsed using `result.Uint()`.

## Modifiers and path chaining

A modifier is a path component that performs custom processing on the YAML.

Multiple paths can be "chained" together using the pipe character. This is useful for getting results from a modified query.

For example, using the built-in `@reverse` modifier on the above YAML document, we'll get the `children` array and reverse the order:

```
"children|@reverse"    >> ["Jack","Alex","Sara"]
"children|@reverse|0"  >> "Jack"
```

There are currently the following built-in modifiers:

- `@reverse`: Reverse an array or the members of an object.
- `@ugly`: Remove all whitespace from a YAML document.
- `@pretty`: Make the YAML document more human readable.
- `@this`: Returns the current element. It can be used to retrieve the root element.
- `@valid`: Ensure the YAML document is valid.
- `@flatten`: Flattens an array.
- `@join`: Joins multiple objects into a single object.
- `@keys`: Returns an array of keys for an object.
- `@values`: Returns an array of values for an object.

### Modifier arguments

A modifier may accept an optional argument. The argument can be a valid YAML document or just characters.

For example, the `@pretty` modifier accepts options that can customize the output.

### Custom modifiers

You can also add custom modifiers.

```go
gyaml.AddModifier("case", func(yaml, arg string) string {
  if arg == "upper" {
    return strings.ToUpper(yaml)
  }
  if arg == "lower" {
    return strings.ToLower(yaml)
  }
  return yaml
})
```

```
"children.0|@case:upper"   >> "SARA"
```

## Get nested array values

Suppose you want all the last names from the following YAML:

```yaml
programmers:
  - firstName: Janet
    lastName: McLaughlin
  - firstName: Elliotte
    lastName: Hunter
  - firstName: Jason
    lastName: Harold
```

You would use the path `programmers.#.lastName` like such:

```go
result := gyaml.Get(yaml, "programmers.#.lastName")
for _, name := range result.Array() {
    println(name.String())
}
```

You can also query an object inside an array:

```go
name := gyaml.Get(yaml, `programmers.#(lastName=="Hunter").firstName`)
println(name.String()) // prints "Elliotte"
```

## Iterate through an object or array

The `ForEach` function allows for quickly iterating through an object or array. The key and value are passed to the iterator function for objects. Only the value is passed for arrays. Returning `false` from an iterator will stop iteration.

```go
result := gyaml.Get(yaml, "programmers")
result.ForEach(func(key, value gyaml.Result) bool {
    println(value.String()) 
    return true // keep iterating
})
```

## Simple Parse and Get

There's a `Parse(yaml)` function that will do a simple parse, and `result.Get(path)` that will search a result.

For example, all of these will return the same result:

```go
gyaml.Parse(yaml).Get("name").Get("last")
gyaml.Get(yaml, "name").Get("last")
gyaml.Get(yaml, "name.last")
```

## Check for the existence of a value

Sometimes you just want to know if a value exists.

```go
value := gyaml.Get(yaml, "name.last")
if !value.Exists() {
    println("no last name")
} else {
    println(value.String())
}

// Or as one step
if gyaml.Get(yaml, "name.last").Exists() {
    println("has a last name")
}
```

## Validate YAML

The `Get*` and `Parse*` functions expects that the YAML is well-formed. Bad YAML will not panic, but it may return back unexpected results.

If you are consuming YAML from an unpredictable source then you may want to validate prior to using GYAML.

```go
if !gyaml.Valid(yaml) {
    return errors.New("invalid yaml")
}
value := gyaml.Get(yaml, "name.last")
```

## Unmarshal to a map

To unmarshal to a `map[string]interface{}`:

```go
m, ok := gyaml.Parse(yaml).Value().(map[string]interface{})
if !ok {
    // not a map
}
```

## Working with Bytes

If your YAML is contained in a `[]byte` slice, there's the `GetBytes` function. This is preferred over `Get(string(data), path)`.

```go
var yaml []byte = ...
result := gyaml.GetBytes(yaml, path)
```

If you are using the `gyaml.GetBytes(yaml, path)` function and you want to avoid converting `result.Raw` to a `[]byte`, then you can use this pattern:

```go
var yaml []byte = ...
result := gyaml.GetBytes(yaml, path)
var raw []byte
if result.Index > 0 {
    raw = yaml[result.Index:result.Index+len(result.Raw)]
} else {
    raw = []byte(result.Raw)
}
```

This is a best-effort no allocation sub slice of the original YAML. This method utilizes the `result.Index` field, which is the position of the raw data in the original YAML. It's possible that the value of `result.Index` equals zero, in which case the `result.Raw` is converted to a `[]byte`.

## ⚡ Performance

Benchmarks for [gyaml](https://github.com/m4l1c1ou5/gyaml) alongside [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)

```
BenchmarkGYAMLGet-8                       615477              1932 ns/op            2882 B/op         36 allocs/op
BenchmarkGYAMLGetBytes-8                  614497              2001 ns/op            2882 B/op         36 allocs/op
BenchmarkGYAMLGetComplex-8                 76069             15730 ns/op           17504 B/op        267 allocs/op
BenchmarkGYAMLGetArray-8                   35044             34759 ns/op           36509 B/op        569 allocs/op
BenchmarkGYAMLParse-8                   1000000000               0.3381 ns/op          0 B/op          0 allocs/op
BenchmarkGYAMLValid-8                      71376             16656 ns/op           18640 B/op        273 allocs/op
BenchmarkYAMLv3UnmarshalMap-8              70812             17210 ns/op           18632 B/op        273 allocs/op
BenchmarkYAMLv3UnmarshalStruct-8           73639             17369 ns/op           16608 B/op        230 allocs/op
BenchmarkGYAMLMultipleGets-8              198769              6816 ns/op            8648 B/op        110 allocs/op
BenchmarkYAMLv3MultipleGets-8              72682             17468 ns/op           18632 B/op        273 allocs/op
BenchmarkGYAMLDeepPath-8                   27262             44245 ns/op           47480 B/op        748 allocs/op
BenchmarkYAMLv3DeepPath-8                  32686             36935 ns/op           37096 B/op        613 allocs/op
BenchmarkGYAMLQueryAll-8                   32304             35570 ns/op           40040 B/op        575 allocs/op
BenchmarkGYAMLQueryConditional-8           35884             37808 ns/op           37255 B/op        548 allocs/op
BenchmarkGYAMLForEach-8                    39108             30795 ns/op           57048 B/op        402 allocs/op
BenchmarkGYAMLArray-8                      10000            114170 ns/op          226442 B/op       1484 allocs/op
BenchmarkGYAMLMap-8                        44298             27203 ns/op           38856 B/op        445 allocs/op
BenchmarkGYAMLGetString-8                 616413              1894 ns/op            2800 B/op         36 allocs/op
BenchmarkGYAMLGetInt-8                    613545              2004 ns/op            2800 B/op         36 allocs/op
BenchmarkGYAMLGetBool-8                    35438             34354 ns/op           34744 B/op        567 allocs/op
BenchmarkGYAMLLargeDocument-8              28072             42936 ns/op           43584 B/op        680 allocs/op
BenchmarkYAMLv3LargeDocument-8             32971             37913 ns/op           33560 B/op        564 allocs/op
```


System Information on which benchmarks were recorded:
```
Go Version: go version go1.25.1 darwin/arm64
System: macOS
Model Name: MacBook Air
Chip: Apple M3
Memory: 8 GB
```

Last run: Nov 22, 2024

Benchmarking script can be found [here](https://github.com/m4l1c1ou5/gyaml-benchmarks).
