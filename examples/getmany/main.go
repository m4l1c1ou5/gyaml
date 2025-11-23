package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// GetManyExample demonstrates GetMany and GetManyBytes operations
func main() {
	const yaml = `
name:
  first: Tom
  last: Anderson
  middle: James
age: 37
city: New York
active: true
children:
  - Sara
  - Alex
  - Jack
scores:
  math: 95
  science: 88
  english: 92
`

	fmt.Println("=== GetMany Examples ===\n")

	// 1. GetMany - retrieve multiple paths at once
	fmt.Println("1. GetMany - Multiple paths:")
	results := gyaml.GetMany(yaml, "name.first", "age", "city").Array()
	for i, result := range results {
		fmt.Printf("   [%d] %s\n", i, result.String())
	}

	// 2. GetMany with mixed types
	fmt.Println("\n2. GetMany - Mixed types:")
	mixed := gyaml.GetMany(yaml, "name.first", "age", "active", "children.0").Array()
	fmt.Println("   String:", mixed[0].String())
	fmt.Println("   Int:", mixed[1].Int())
	fmt.Println("   Bool:", mixed[2].Bool())
	fmt.Println("   Array element:", mixed[3].String())

	// 3. GetMany with nested paths
	fmt.Println("\n3. GetMany - Nested paths:")
	nested := gyaml.GetMany(yaml, "name.first", "name.last", "name.middle").Array()
	for i, val := range nested {
		fmt.Printf("   Name part %d: %s\n", i+1, val.String())
	}

	// 4. GetMany with array paths
	fmt.Println("\n4. GetMany - Array elements:")
	arrayElems := gyaml.GetMany(yaml, "children.0", "children.1", "children.2").Array()
	for i, child := range arrayElems {
		fmt.Printf("   Child %d: %s\n", i+1, child.String())
	}

	// 5. GetMany with object paths
	fmt.Println("\n5. GetMany - Multiple scores:")
	scoreResults := gyaml.GetMany(yaml, "scores.math", "scores.science", "scores.english").Array()
	for i, score := range scoreResults {
		subjects := []string{"Math", "Science", "English"}
		fmt.Printf("   %s: %d\n", subjects[i], score.Int())
	}

	// 6. GetMany with non-existent paths
	fmt.Println("\n6. GetMany - Including non-existent path:")
	withMissing := gyaml.GetMany(yaml, "name.first", "nonexistent", "age").Array()
	for i, result := range withMissing {
		if result.Exists() {
			fmt.Printf("   [%d] %s (exists)\n", i, result.String())
		} else {
			fmt.Printf("   [%d] <not found>\n", i)
		}
	}

	// 7. GetMany iterating over results
	fmt.Println("\n7. GetMany - Iterate with Array():")
	many := gyaml.GetMany(yaml, "name.first", "age", "children.0")
	for i, v := range many.Array() {
		fmt.Printf("   [%d] %s\n", i, v.String())
	}

	// 8. GetManyBytes with byte slice
	fmt.Println("\n8. GetManyBytes - Working with bytes:")
	yamlBytes := []byte(yaml)
	byteResults := gyaml.GetManyBytes(yamlBytes, "name.last", "city", "active").Array()
	for i, result := range byteResults {
		fmt.Printf("   [%d] %s\n", i, result.String())
	}

	// 9. GetMany for building structs
	fmt.Println("\n9. GetMany - Building a struct-like result:")
	person := gyaml.GetMany(yaml, "name.first", "name.last", "age", "city").Array()
	fmt.Printf("   Person: %s %s, Age: %d, City: %s\n",
		person[0].String(),
		person[1].String(),
		person[2].Int(),
		person[3].String())

	// 10. GetMany vs multiple Get calls
	fmt.Println("\n10. GetMany is more efficient than multiple Get calls:")
	fmt.Println("   // Instead of:")
	fmt.Println("   // first := gyaml.Get(yaml, \"name.first\")")
	fmt.Println("   // last := gyaml.Get(yaml, \"name.last\")")
	fmt.Println("   // age := gyaml.Get(yaml, \"age\")")
	fmt.Println("   // Use:")
	fmt.Println("   // results := gyaml.GetMany(yaml, \"name.first\", \"name.last\", \"age\").Array()")

	efficient := gyaml.GetMany(yaml, "name.first", "name.last", "age").Array()
	fmt.Printf("   Results: %s %s, %d\n",
		efficient[0].String(),
		efficient[1].String(),
		efficient[2].Int())

	// 11. GetMany with queries
	fmt.Println("\n11. GetMany - With queries:")
	const friendsYaml = `
friends:
  - first: Dale
    last: Murphy
  - first: Roger
    last: Craig
`
	queryResults := gyaml.GetMany(friendsYaml,
		"friends.0.first",
		"friends.1.first",
		"friends.#.last").Array()

	if len(queryResults) >= 3 {
		fmt.Println("   Friend 1 first name:", queryResults[0].String())
		fmt.Println("   Friend 2 first name:", queryResults[1].String())
		fmt.Print("   All last names: ")
		for _, name := range queryResults[2].Array() {
			fmt.Print(name.String(), " ")
		}
		fmt.Println()
	} else {
		fmt.Println("   Note: Some query paths did not return results")
	}
}
