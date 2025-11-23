package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// ForEachExample demonstrates ForEach iteration
func ForEachExample() {
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
friends:
  - first: Dale
    last: Murphy
    age: 44
  - first: Roger
    last: Craig
    age: 68
scores:
  math: 95
  science: 88
  english: 92
`

	fmt.Println("=== ForEach Examples ===\n")

	// 1. ForEach on array - only value is used
	fmt.Println("1. ForEach on array (children):")
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		fmt.Println("   Child:", value.String())
		return true // continue iteration
	})

	// 2. ForEach on object - both key and value are used
	fmt.Println("\n2. ForEach on object (name):")
	gyaml.Get(yaml, "name").ForEach(func(key, value gyaml.Result) bool {
		fmt.Printf("   %s: %s\n", key.String(), value.String())
		return true
	})

	// 3. Early termination - stop after first item
	fmt.Println("\n3. Early termination (stop after first child):")
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		fmt.Println("   First child:", value.String())
		return false // stop iteration
	})

	// 4. ForEach with conditional processing
	fmt.Println("\n4. ForEach with condition (scores >= 90):")
	gyaml.Get(yaml, "scores").ForEach(func(key, value gyaml.Result) bool {
		if value.Int() >= 90 {
			fmt.Printf("   High score in %s: %d\n", key.String(), value.Int())
		}
		return true
	})

	// 5. ForEach on array of objects
	fmt.Println("\n5. ForEach on array of objects (friends):")
	gyaml.Get(yaml, "friends").ForEach(func(key, value gyaml.Result) bool {
		firstName := value.Get("first").String()
		lastName := value.Get("last").String()
		age := value.Get("age").Int()
		fmt.Printf("   Friend: %s %s, Age: %d\n", firstName, lastName, age)
		return true
	})

	// 6. Nested ForEach
	fmt.Println("\n6. Nested ForEach (friends and their properties):")
	gyaml.Get(yaml, "friends").ForEach(func(idx, friend gyaml.Result) bool {
		fmt.Printf("   Friend %d:\n", idx.Int())
		friend.ForEach(func(key, value gyaml.Result) bool {
			fmt.Printf("     %s: %s\n", key.String(), value.String())
			return true
		})
		return true
	})

	// 7. Counting items with ForEach
	fmt.Println("\n7. Counting with ForEach:")
	count := 0
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		count++
		return true
	})
	fmt.Printf("   Total children: %d\n", count)

	// 8. Collecting values with ForEach
	fmt.Println("\n8. Collecting values:")
	var names []string
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		names = append(names, value.String())
		return true
	})
	fmt.Printf("   Collected names: %v\n", names)

	// 9. ForEach with key index for arrays
	fmt.Println("\n9. ForEach with index on array:")
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		// For arrays, key is the index
		fmt.Printf("   Index %d: %s\n", key.Int(), value.String())
		return true
	})

	// 10. Conditional early exit
	fmt.Println("\n10. Stop iteration when condition met:")
	gyaml.Get(yaml, "children").ForEach(func(key, value gyaml.Result) bool {
		fmt.Println("   Processing:", value.String())
		if value.String() == "Alex" {
			fmt.Println("   Found Alex! Stopping.")
			return false
		}
		return true
	})

	// 11. ForEach on empty result
	fmt.Println("\n11. ForEach on non-existent path:")
	called := false
	gyaml.Get(yaml, "nonexistent").ForEach(func(key, value gyaml.Result) bool {
		called = true
		return true
	})
	fmt.Printf("   Callback called: %v\n", called)

	// 12. Building a summary with ForEach
	fmt.Println("\n12. Building summary from object:")
	summary := ""
	gyaml.Get(yaml, "scores").ForEach(func(key, value gyaml.Result) bool {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%s: %d", key.String(), value.Int())
		return true
	})
	fmt.Printf("   Summary: %s\n", summary)
}

func main() {
	ForEachExample()
}
