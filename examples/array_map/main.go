package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// ArrayMapExample demonstrates Array() and Map() operations
func ArrayMapExample() {
	const yaml = `
name:
  first: Tom
  last: Anderson
age: 37
children:
  - Sara
  - Alex
  - Jack
scores:
  math: 95
  science: 88
friends:
  - first: Dale
    last: Murphy
  - first: Roger
    last: Craig
`

	fmt.Println("=== Array and Map Examples ===\n")

	// 1. Array() on simple array
	fmt.Println("1. Array() on simple array:")
	children := gyaml.Get(yaml, "children")
	for i, child := range children.Array() {
		fmt.Printf("   [%d] %s\n", i, child.String())
	}

	// 2. Map() on object
	fmt.Println("\n2. Map() on object:")
	nameMap := gyaml.Get(yaml, "name").Map()
	for key, value := range nameMap {
		fmt.Printf("   %s: %s\n", key, value.String())
	}

	// 3. Array() with objects
	fmt.Println("\n3. Array() of objects:")
	friendsArray := gyaml.Get(yaml, "friends").Array()
	for i, friend := range friendsArray {
		fmt.Printf("   Friend %d: %s %s\n", i+1,
			friend.Get("first").String(),
			friend.Get("last").String())
	}

	// 4. Combining Array() and Map()
	fmt.Println("\n4. Combining Array() and Map():")
	for i, friend := range gyaml.Get(yaml, "friends").Array() {
		fmt.Printf("   Friend %d:\n", i+1)
		friendMap := friend.Map()
		for k, v := range friendMap {
			fmt.Printf("     %s: %s\n", k, v.String())
		}
	}

	// 5. Converting to native Go types
	fmt.Println("\n5. Convert to native Go types:")
	childrenValue := gyaml.Get(yaml, "children").Value()
	if arr, ok := childrenValue.([]interface{}); ok {
		fmt.Printf("   Children: %v\n", arr)
	}
}

func main() {
	ArrayMapExample()
}
