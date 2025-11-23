package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// PathSyntaxExample demonstrates various path syntax features
func PathSyntaxExample() {
	const yaml = `
name:
  first: Tom
  last: Anderson
age: 37
children:
  - Sara
  - Alex
  - Jack
fav.movie: Deer Hunter
favorites:
  color: blue
  food: pizza
friends:
  - first: Dale
    last: Murphy
    age: 44
  - first: Roger
    last: Craig
    age: 68
`

	fmt.Println("=== Path Syntax Examples ===\n")

	// 1. Simple dot notation
	fmt.Println("1. Simple path (name.last):")
	fmt.Println("  ", gyaml.Get(yaml, "name.last").String())

	// 2. Array index access
	fmt.Println("\n2. Array index (children.1):")
	fmt.Println("  ", gyaml.Get(yaml, "children.1").String())

	// 3. Array count
	fmt.Println("\n3. Array count (children.#):")
	fmt.Println("  ", gyaml.Get(yaml, "children.#").Int())

	// 4. Wildcard * - matches any characters
	fmt.Println("\n4. Wildcard * (child*.2):")
	fmt.Println("  ", gyaml.Get(yaml, "child*.2").String())

	// 5. Wildcard ? - matches single character
	fmt.Println("\n5. Wildcard ? (c?ildren.0):")
	fmt.Println("  ", gyaml.Get(yaml, "c?ildren.0").String())

	// 6. Escaped dot for keys containing dots
	fmt.Println("\n6. Escaped dot (fav\\.movie):")
	fmt.Println("  ", gyaml.Get(yaml, `fav\.movie`).String())

	// 7. Get all elements from nested arrays
	fmt.Println("\n7. Nested array path (friends.#.first):")
	result := gyaml.Get(yaml, "friends.#.first")
	for _, name := range result.Array() {
		fmt.Println("  ", name.String())
	}

	// 8. Get nested value from array element
	fmt.Println("\n8. Array element nested value (friends.1.last):")
	fmt.Println("  ", gyaml.Get(yaml, "friends.1.last").String())

	// 9. Multiple wildcards
	fmt.Println("\n9. Multiple wildcards (fav*.c*r):")
	fmt.Println("  ", gyaml.Get(yaml, "fav*.c*r").String())

	// 10. Get entire array
	fmt.Println("\n10. Get entire array (children):")
	children := gyaml.Get(yaml, "children")
	fmt.Println("   ", children.Raw)

	// 11. First element of array
	fmt.Println("\n11. First element (children.0):")
	fmt.Println("   ", gyaml.Get(yaml, "children.0").String())

	// 12. Last element of array (using count)
	fmt.Println("\n12. Using wildcard with nested object (fav*.food):")
	fmt.Println("   ", gyaml.Get(yaml, "fav*.food").String())
}

func main() {
	PathSyntaxExample()
}
