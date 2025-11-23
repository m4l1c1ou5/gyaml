# gyaml Examples

This directory contains comprehensive examples demonstrating all the features of the GYAML library.

## Directory Structure

Each example is in its own subdirectory with a `main.go` file. This organization prevents conflicts from having multiple main packages.

```
examples/
├── basic/                  - Basic operations
├── path_syntax/            - Path syntax features
├── queries/                - Query operations
├── modifiers/              - Built-in modifiers
├── custom_modifiers/       - Custom modifiers
├── result_type/            - Result type and methods
├── foreach/                - ForEach iteration
├── validation/             - YAML validation
├── bytes/                  - Working with bytes
├── array_map/              - Array and Map operations
├── getmany/                - GetMany operations
└── README.md
```

## Running Examples

Run the example by:

```bash
go run <example_directory>/main.go
```

## Available Examples

### 1. **basic/** - Basic Operations
Demonstrates fundamental GYAML operations:
- `Get()` - Retrieve values from YAML
- `Parse()` - Parse YAML documents
- Type conversions: `String()`, `Int()`, `Float()`, `Bool()`
- Chaining `Get()` calls
- Accessing `Raw` and `Type` properties

### 2. **path_syntax/** - Path Syntax
Shows various path syntax features:
- Simple dot notation (`name.last`)
- Array indexing (`children.1`)
- Array counting (`children.#`)
- Wildcards (`*` and `?`)
- Escaped dots for keys with dots (`fav\.movie`)
- Nested paths (`friends.#.first`)

### 3. **queries/** - Query Operations
Demonstrates query syntax and operators:
- Equality: `==` and `!=`
- Comparison: `<`, `<=`, `>`, `>=`
- Pattern matching: `%` (like) and `!%` (not like)
- First match: `#(...)`
- All matches: `#(...)#`
- Nested array queries

### 4. **modifiers/** - Built-in Modifiers
Shows all built-in modifiers:
- `@reverse` - Reverse arrays/objects
- `@keys` - Get object keys
- `@values` - Get object values
- `@this` - Get current element
- `@valid` - Validate YAML
- `@flatten` - Flatten nested arrays
- `@ugly` - Remove whitespace
- `@pretty` - Pretty print
- `@join` - Join objects
- Chaining modifiers

### 5. **custom_modifiers/** - Custom Modifiers
Learn how to create custom modifiers:
- `AddModifier()` function
- Modifier with arguments
- Case conversion example
- String manipulation examples
- Chaining custom modifiers

### 6. **result_type/** - Result Type
Comprehensive coverage of the Result type:
- Properties: `Type`, `Str`, `Num`, `Raw`, `Index`, `Indexes`
- Methods: `Exists()`, `Value()`, `String()`, `Int()`, `Uint()`, `Float()`, `Bool()`, `Time()`
- `Array()` and `Map()` methods
- `Get()` on Result objects
- Type checking: `IsArray()`, `IsObject()`
- Handling non-existent values

### 7. **foreach/** - ForEach Iteration
Demonstrates iteration patterns:
- Iterating over arrays
- Iterating over objects
- Early termination
- Conditional processing
- Nested iteration
- Collecting values
- Counting with ForEach

### 8. **validation/** - YAML Validation
Shows validation techniques:
- `Valid()` function
- Valid vs invalid YAML examples
- Best practices for validation
- Using `@valid` modifier
- Validating untrusted input

### 9. **bytes/** - Working with Bytes
Demonstrates byte slice operations:
- `GetBytes()` - Preferred over `Get(string(data), path)`
- `GetManyBytes()` - Multiple paths
- Efficient byte extraction using `Index` field
- Zero-allocation patterns
- Performance considerations

### 10. **array_map/** - Array and Map Operations
Shows array and map manipulation:
- `Array()` method
- `Map()` method
- Iteration patterns
- Type conversion to native Go types
- Combining operations
- Edge cases

### 11. **getmany/** - GetMany Operations
Demonstrates retrieving multiple paths efficiently:
- `GetMany()` - Multiple paths at once
- `GetManyBytes()` - Byte slice variant
- Mixed type handling
- Building struct-like results
- Performance benefits over multiple Get calls

## Common Patterns

### Safe Value Extraction
```go
value := gyaml.Get(yaml, "path.to.value")
if value.Exists() {
    // Use value.String(), value.Int(), etc.
}
```

### Iterating Results
```go
results := gyaml.Get(yaml, "items.#.name")
for _, item := range results.Array() {
    fmt.Println(item.String())
}
```

### Validating Before Use
```go
if !gyaml.Valid(yaml) {
    return errors.New("invalid yaml")
}
result := gyaml.Get(yaml, "path")
```

### Working with Bytes
```go
var data []byte = ...
result := gyaml.GetBytes(data, "path")
```

## More Information

- [Main README](../README.md) - Full documentation
- [SYNTAX.md](../SYNTAX.md) - Detailed path syntax reference
- [GitHub Repository](https://github.com/m4l1c1ou5/gyaml)

## Contributing Examples

If you have additional examples that demonstrate useful patterns, please contribute them to this directory!
