package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// ModifiersExample demonstrates all built-in modifiers
func ModifiersExample() {
	const yaml = `
name:
  first: Tom
  last: Anderson
  middle: James
age: 37
children:
  - Sara
  - Alex
  - Jack
numbers:
  - - 1
    - 2
  - - 3
    - 4
  - - 5
    - 6
person1:
  name: John
  age: 30
person2:
  city: NYC
  country: USA
`

	fmt.Println("=== Modifiers Examples ===\n")

	// 1. @reverse - Reverse an array
	fmt.Println("1. @reverse modifier on array (children|@reverse):")
	reversed := gyaml.Get(yaml, "children|@reverse")
	for _, r := range reversed.Array() {
		fmt.Println("  ", r.String())
	}

	// 2. @reverse on chained path
	fmt.Println("\n2. @reverse with chaining (children|@reverse|0):")
	fmt.Println("  ", gyaml.Get(yaml, "children|@reverse|0").String())

	// 3. @keys - Get object keys
	fmt.Println("\n3. @keys modifier (name|@keys):")
	keys := gyaml.Get(yaml, "name|@keys")
	for _, k := range keys.Array() {
		fmt.Println("  ", k.String())
	}

	// 4. @values - Get object values
	fmt.Println("\n4. @values modifier (name|@values):")
	values := gyaml.Get(yaml, "name|@values")
	for _, v := range values.Array() {
		fmt.Println("  ", v.String())
	}

	// 5. @this - Get current element
	fmt.Println("\n5. @this modifier (returns root):")
	thisResult := gyaml.Get(yaml, "@this")
	fmt.Println("   Type:", thisResult.Type)
	fmt.Println("   Has content:", len(thisResult.Raw) > 0)

	// 6. @valid - Validate YAML
	fmt.Println("\n6. @valid modifier:")
	valid := gyaml.Get(yaml, "@valid")
	fmt.Println("  ", valid.Raw)

	// 7. @flatten - Flatten nested arrays
	fmt.Println("\n7. @flatten modifier (numbers|@flatten):")
	flattened := gyaml.Get(yaml, "numbers|@flatten")
	for _, n := range flattened.Array() {
		fmt.Println("  ", n.String())
	}

	// 8. @ugly - Remove whitespace
	fmt.Println("\n8. @ugly modifier (name|@ugly):")
	ugly := gyaml.Get(yaml, "name|@ugly")
	fmt.Println("  ", ugly.Raw)

	// 9. @pretty - Pretty print
	fmt.Println("\n9. @pretty modifier (name|@pretty):")
	pretty := gyaml.Get(yaml, "name|@pretty")
	fmt.Println(pretty.Raw)

	// 10. @join - Join multiple objects
	fmt.Println("\n10. @join modifier:")
	const joinYaml = `
objects:
  - name: John
    age: 30
  - city: NYC
    country: USA
`
	joined := gyaml.Get(joinYaml, "objects|@join")
	fmt.Println(joined.Raw)

	// 11. Multiple modifiers chained
	fmt.Println("\n11. Chain multiple modifiers (children|@reverse|@ugly):")
	chained := gyaml.Get(yaml, "children|@reverse|@ugly")
	fmt.Println("  ", chained.Raw)

	// 12. Modifier with path continuation
	fmt.Println("\n12. Modifier with continued path (name|@keys|0):")
	fmt.Println("  ", gyaml.Get(yaml, "name|@keys|0").String())
}

func main() {
	ModifiersExample()
}
