# gyaml Path Syntax

This document describes the gyaml path syntax.

## Basic Path

A path is a series of keys seperated by a dot.

The dot `.` character separates path segments.

For example, given this YAML:

```yaml
name:
  first: Janet
  last: Prichard
```

The path `name.first` would return `Janet`.

## Arrays

Arrays can be accessed by index using the index as the key.

For example, given this YAML:

```yaml
friends:
  - Dale
  - Roger
  - Jane
```

- `friends.0` returns `Dale`
- `friends.1` returns `Roger`
- `friends.2` returns `Jane`

## Array Length

The `#` character returns the number of elements in an array.

```yaml
friends:
  - Dale
  - Roger
  - Jane
```

- `friends.#` returns `3`

## Nested Arrays

You can access nested values within arrays:

```yaml
friends:
  - name: Dale
    age: 44
  - name: Roger
    age: 68
  - name: Jane
    age: 47
```

- `friends.0.name` returns `Dale`
- `friends.1.age` returns `68`
- `friends.#.name` returns `["Dale","Roger","Jane"]` (all names)

## Wildcards

Wildcard characters `*` and `?` can be used in keys:

- `*` matches zero or more characters
- `?` matches exactly one character

```yaml
child_1: first
child_2: second
child_3: third
```

- `child_*` matches all three keys
- `child_?` matches `child_1` only if single char after underscore

## Escape Character

Use backslash `\` to escape special characters:

```yaml
fav.movie: Deer Hunter
```

- `fav\.movie` returns `Deer Hunter` (dot is escaped)

## Queries

You can query arrays using `#(...)` for the first match, or `#(...)#` for all matches.

### Comparison Operators

Queries support these comparison operators:

- `==` equal
- `!=` not equal
- `<` less than
- `<=` less than or equal
- `>` greater than
- `>=` greater than or equal

### Pattern Matching

- `%` like (wildcard pattern match)
- `!%` not like

### Query Examples

Given this YAML:

```yaml
friends:
  - name: Dale
    age: 44
  - name: Roger
    age: 68
  - name: Jane
    age: 47
```

- `friends.#(name=="Dale")` returns the first friend named Dale
- `friends.#(name=="Dale").age` returns `44`
- `friends.#(age>45)#.name` returns `["Roger","Jane"]`
- `friends.#(name%"D*").age` returns `44` (Dale matches pattern D*)

### Nested Queries

Queries can be nested:

```yaml
friends:
  - name: Dale
    pets:
      - dog
      - cat
  - name: Roger
    pets:
      - bird
```

- `friends.#(pets.#(=="dog"))` returns Dale's object (has a dog)

## Modifiers

Modifiers are special functions that transform the result. They start with `@`.

### Built-in Modifiers

- `@reverse` - Reverse an array or object members
- `@ugly` - Remove whitespace from YAML  
- `@pretty` - Format YAML for readability
- `@this` - Return the current element (useful for root)
- `@valid` - Return `true` if YAML is valid, `false` otherwise
- `@flatten` - Flatten nested arrays
- `@join` - Join multiple objects into one
- `@keys` - Return array of object keys
- `@values` - Return array of object values

### Modifier Examples

```yaml
children:
  - Sara
  - Alex
  - Jack
```

- `children|@reverse` returns `["Jack","Alex","Sara"]`
- `children|@reverse.0` returns `Jack`

### Modifier Arguments

Some modifiers accept arguments after a colon:

```yaml
name: john
```

- `name|@custom:upper` (if custom modifier is defined)

### Chaining

Multiple modifiers can be chained with `|`:

```yaml
numbers:
  - 1
  - 2
  - 3
```

- `numbers|@reverse|@this` returns reversed array

## Multipaths

Get multiple paths at once with `GetMany`:

```go
result := gyaml.GetMany(yaml, "name.first", "name.last", "age")
// Returns array: ["Janet", "Prichard", 47]
```

## YAML Lines

Lines of YAML can be treated as an array using the `..` prefix:

```yaml
name: Gilbert
age: 61
---
name: Alexa
age: 34
```

Or JSON-lines style:

```
{"name": "Gilbert", "age": 61}
{"name": "Alexa", "age": 34}
```

- `..#` returns `2` (number of lines)
- `..0` returns first line
- `..#.name` returns `["Gilbert","Alexa"]`

## Dotted Keys

Keys that contain dots must be escaped:

```yaml
fav.movie: Deer Hunter
```

Use `fav\.movie` to access this value.

## Complex Example

```yaml
users:
  - name: Tom
    age: 37
    active: true
    tags:
      - admin
      - developer
  - name: Jane
    age: 29
    active: true
    tags:
      - user
  - name: Bob
    age: 45
    active: false
    tags:
      - developer
```

Queries:

- `users.#(active==true)#.name` → `["Tom","Jane"]`
- `users.#(age>30)#.name` → `["Tom","Bob"]`
- `users.#(tags.#(=="admin")).name` → `Tom`
- `users.#.tags|@flatten` → `["admin","developer","user","developer"]`
- `users.#.age` → `[37,29,45]`

## Performance Tips

1. **Be specific**: More specific paths are faster
2. **Avoid wildcards when possible**: Exact matches are faster than wildcards
3. **Use queries efficiently**: `#(...)` stops at first match, `#(...)#` checks all
4. **Cache parsed results**: Use `Parse()` if making multiple queries on same YAML

## Edge Cases

### Empty Path

An empty path `""` returns the entire YAML document.

### Non-existent Paths

Non-existent paths return a `Null` type result where `Exists()` returns `false`.

### Type Mismatches

Accessing an array index on a non-array returns `Null`.
Accessing an object key on a non-object returns `Null`.

### Special Characters in Keys

Keys with special characters should be escaped:

```yaml
my-key: value
my.key: value
my*key: value
```

Access with:
- `my-key` (no escape needed for dash)
- `my\.key` (escape dot)
- `my\*key` (escape asterisk)
