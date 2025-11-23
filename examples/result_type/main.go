package main

import (
	"fmt"
	"time"

	"github.com/m4l1c1ou5/gyaml"
)

// ResultTypeExample demonstrates Result type and its methods
func ResultTypeExample() {
	const yaml = `
name: John Doe
age: 30
salary: 75000.50
active: true
nullValue: null
timestamp: 2024-01-15T10:30:00Z
tags:
  - golang
  - yaml
  - programming
address:
  street: 123 Main St
  city: New York
  zip: 10001
bigInt: 9007199254740993
`

	fmt.Println("=== Result Type Examples ===\n")

	// 1. Type property
	fmt.Println("1. Result.Type property:")
	fmt.Println("   String type:", gyaml.Get(yaml, "name").Type)
	fmt.Println("   Number type:", gyaml.Get(yaml, "age").Type)
	fmt.Println("   Bool type:", gyaml.Get(yaml, "active").Type)
	fmt.Println("   Null type:", gyaml.Get(yaml, "nullValue").Type)
	fmt.Println("   Array type:", gyaml.Get(yaml, "tags").Type)
	fmt.Println("   Object type:", gyaml.Get(yaml, "address").Type)

	// 2. Str property
	fmt.Println("\n2. Result.Str property:")
	fmt.Println("  ", gyaml.Get(yaml, "name").Str)

	// 3. Num property
	fmt.Println("\n3. Result.Num property:")
	fmt.Println("  ", gyaml.Get(yaml, "age").Num)

	// 4. Raw property
	fmt.Println("\n4. Result.Raw property:")
	fmt.Println("   String:", gyaml.Get(yaml, "name").Raw)
	fmt.Println("   Array:", gyaml.Get(yaml, "tags").Raw)
	fmt.Println("   Object:", gyaml.Get(yaml, "address").Raw)

	// 5. Index property
	fmt.Println("\n5. Result.Index property:")
	result := gyaml.Get(yaml, "name")
	fmt.Println("   Index:", result.Index)

	// 6. Exists() method
	fmt.Println("\n6. Result.Exists() method:")
	fmt.Println("   name exists:", gyaml.Get(yaml, "name").Exists())
	fmt.Println("   nonexistent exists:", gyaml.Get(yaml, "nonexistent").Exists())

	// 7. Value() method
	fmt.Println("\n7. Result.Value() interface{}:")
	value := gyaml.Get(yaml, "name").Value()
	fmt.Printf("   Type: %T, Value: %v\n", value, value)

	// 8. String() method
	fmt.Println("\n8. Result.String() method:")
	fmt.Println("  ", gyaml.Get(yaml, "name").String())

	// 9. Int() method
	fmt.Println("\n9. Result.Int() method:")
	fmt.Println("  ", gyaml.Get(yaml, "age").Int())

	// 10. Uint() method for large integers
	fmt.Println("\n10. Result.Uint() for 64-bit integers:")
	bigIntResult := gyaml.Get(yaml, "bigInt")
	fmt.Println("   Uint:", bigIntResult.Uint())
	fmt.Println("   Note: This number cannot be precisely represented as float64")

	// 11. Float() method
	fmt.Println("\n11. Result.Float() method:")
	fmt.Println("  ", gyaml.Get(yaml, "salary").Float())

	// 12. Bool() method
	fmt.Println("\n12. Result.Bool() method:")
	fmt.Println("  ", gyaml.Get(yaml, "active").Bool())

	// 13. Time() method
	fmt.Println("\n13. Result.Time() method:")
	timeResult := gyaml.Get(yaml, "timestamp")
	parsedTime := timeResult.Time()
	fmt.Println("  ", parsedTime.Format(time.RFC3339))

	// 14. Array() method
	fmt.Println("\n14. Result.Array() method:")
	tags := gyaml.Get(yaml, "tags")
	for i, tag := range tags.Array() {
		fmt.Printf("   [%d] %s\n", i, tag.String())
	}

	// 15. Map() method
	fmt.Println("\n15. Result.Map() method:")
	addressMap := gyaml.Get(yaml, "address").Map()
	for key, val := range addressMap {
		fmt.Printf("   %s: %s\n", key, val.String())
	}

	// 16. Get() method on Result
	fmt.Println("\n16. Result.Get() method:")
	address := gyaml.Get(yaml, "address")
	fmt.Println("   City:", address.Get("city").String())

	// 17. IsArray() and IsObject() helpers
	fmt.Println("\n17. Type checking:")
	fmt.Println("   tags is array:", gyaml.Get(yaml, "tags").IsArray())
	fmt.Println("   address is object:", gyaml.Get(yaml, "address").IsObject())
	fmt.Println("   name is array:", gyaml.Get(yaml, "name").IsArray())

	// 18. Indexes property (for query results)
	fmt.Println("\n18. Result.Indexes (with array query):")
	const arrayYaml = `
items:
  - id: 1
    name: Item1
  - id: 2
    name: Item2
  - id: 3
    name: Item3
`
	queryResult := gyaml.Get(arrayYaml, "items.#.name")
	fmt.Println("   Indexes:", queryResult.Indexes)

	// 19. Handling non-existent values
	fmt.Println("\n19. Non-existent value handling:")
	missing := gyaml.Get(yaml, "nonexistent.path")
	fmt.Println("   Exists:", missing.Exists())
	fmt.Println("   String (default empty):", missing.String())
	fmt.Println("   Int (default 0):", missing.Int())
	fmt.Println("   Bool (default false):", missing.Bool())
}

func main() {
	ResultTypeExample()
}
