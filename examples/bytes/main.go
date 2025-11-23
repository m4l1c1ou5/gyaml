package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// BytesExample demonstrates working with byte slices
func BytesExample() {
	yamlBytes := []byte(`
name:
  first: Tom
  last: Anderson
age: 37
active: true
children:
  - Sara
  - Alex
  - Jack
address:
  street: 123 Main St
  city: New York
  zip: 10001
`)

	fmt.Println("=== Working with Bytes Examples ===\n")

	// 1. GetBytes - preferred over Get(string(data), path)
	fmt.Println("1. GetBytes basic usage:")
	result := gyaml.GetBytes(yamlBytes, "name.last")
	fmt.Println("   Name:", result.String())

	// 2. GetBytes with nested path
	fmt.Println("\n2. GetBytes with nested path:")
	city := gyaml.GetBytes(yamlBytes, "address.city")
	fmt.Println("   City:", city.String())

	// 3. GetBytes with array
	fmt.Println("\n3. GetBytes with array:")
	children := gyaml.GetBytes(yamlBytes, "children")
	for _, child := range children.Array() {
		fmt.Println("   Child:", child.String())
	}

	// 4. GetBytes with array index
	fmt.Println("\n4. GetBytes with array index:")
	firstChild := gyaml.GetBytes(yamlBytes, "children.0")
	fmt.Println("   First child:", firstChild.String())

	// 5. Using result.Index for efficient byte slice extraction
	fmt.Println("\n5. Efficient byte extraction using Index:")
	nameResult := gyaml.GetBytes(yamlBytes, "name")
	var raw []byte
	if nameResult.Index > 0 {
		raw = yamlBytes[nameResult.Index : nameResult.Index+len(nameResult.Raw)]
		fmt.Println("   Zero-allocation slice:", string(raw))
	} else {
		raw = []byte(nameResult.Raw)
		fmt.Println("   Converted to bytes:", string(raw))
	}

	// 6. GetBytes with query
	fmt.Println("\n6. GetBytes with query:")
	yamlWithArray := []byte(`
users:
  - name: Alice
    age: 30
  - name: Bob
    age: 25
  - name: Charlie
    age: 35
`)
	olderUser := gyaml.GetBytes(yamlWithArray, "users.#(age>30).name")
	fmt.Println("   User over 30:", olderUser.String())

	// 7. GetBytes with modifier
	fmt.Println("\n7. GetBytes with modifier:")
	reversed := gyaml.GetBytes(yamlBytes, "children|@reverse")
	for _, child := range reversed.Array() {
		fmt.Println("   Reversed child:", child.String())
	}

	// 8. GetBytes for numbers
	fmt.Println("\n8. GetBytes for numbers:")
	age := gyaml.GetBytes(yamlBytes, "age")
	fmt.Println("   Age (Int):", age.Int())
	fmt.Println("   Age (Float):", age.Float())

	// 9. GetBytes for boolean
	fmt.Println("\n9. GetBytes for boolean:")
	active := gyaml.GetBytes(yamlBytes, "active")
	fmt.Println("   Active:", active.Bool())

	// 10. Efficient raw extraction pattern
	fmt.Println("\n10. Best-effort no-allocation pattern:")
	addressResult := gyaml.GetBytes(yamlBytes, "address")
	var addressRaw []byte
	if addressResult.Index > 0 {
		// Zero-allocation sub-slice
		addressRaw = yamlBytes[addressResult.Index : addressResult.Index+len(addressResult.Raw)]
		fmt.Println("   Method: zero-allocation sub-slice")
	} else {
		// Fallback conversion
		addressRaw = []byte(addressResult.Raw)
		fmt.Println("   Method: fallback conversion")
	}
	fmt.Println("   Address raw:", string(addressRaw))

	// 11. GetManyBytes
	fmt.Println("\n11. GetManyBytes for multiple paths:")
	results := gyaml.GetManyBytes(yamlBytes, "name.first", "age", "children.1")
	for i, r := range results.Array() {
		fmt.Printf("   [%d] %s\n", i, r.String())
	}

	// 12. Comparing GetBytes vs Get
	fmt.Println("\n12. Performance note:")
	fmt.Println("   GetBytes is preferred over Get(string(data), path)")
	fmt.Println("   Reason: Avoids string conversion overhead")

	// GetBytes usage
	result1 := gyaml.GetBytes(yamlBytes, "name.first")
	fmt.Println("   GetBytes result:", result1.String())

	// Get with conversion (less efficient)
	result2 := gyaml.Get(string(yamlBytes), "name.first")
	fmt.Println("   Get result:", result2.String())

	// 13. Using Index field
	fmt.Println("\n13. Understanding Index field:")
	result = gyaml.GetBytes(yamlBytes, "age")
	fmt.Printf("   Index: %d (position in original byte slice)\n", result.Index)
	fmt.Printf("   Raw: %s (the raw YAML value)\n", result.Raw)
	if result.Index > 0 {
		fmt.Printf("   Original bytes at index: %s\n", string(yamlBytes[result.Index:result.Index+len(result.Raw)]))
	}

	// 14. Working with large byte slices
	fmt.Println("\n14. Large byte slice handling:")
	largeYaml := make([]byte, 0, 1024)
	largeYaml = append(largeYaml, []byte(`
data:
  items:
`)...)
	for i := 0; i < 10; i++ {
		largeYaml = append(largeYaml, []byte(fmt.Sprintf("    - item%d\n", i))...)
	}

	itemCount := gyaml.GetBytes(largeYaml, "data.items.#")
	fmt.Printf("   Items count: %d\n", itemCount.Int())
}

func main() {
	BytesExample()
}
