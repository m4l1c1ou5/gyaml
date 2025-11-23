package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// BasicExample demonstrates basic Get and Parse operations
func BasicExample() {
	const yaml = `
name:
  first: Janet
  last: Prichard
age: 47
city: Los Angeles
active: true
`

	fmt.Println("=== Basic Operations ===\n")

	// 1. Simple Get
	fmt.Println("1. Get name.last:")
	value := gyaml.Get(yaml, "name.last")
	fmt.Println("  ", value.String())

	// 2. Get number
	fmt.Println("\n2. Get age:")
	age := gyaml.Get(yaml, "age")
	fmt.Println("   Int:", age.Int())
	fmt.Println("   Float:", age.Float())
	fmt.Println("   Uint:", age.Uint())

	// 3. Get boolean
	fmt.Println("\n3. Get active:")
	active := gyaml.Get(yaml, "active")
	fmt.Println("  ", active.Bool())

	// 4. Parse and Get chaining
	fmt.Println("\n4. Parse and Get chaining:")
	result := gyaml.Parse(yaml).Get("name").Get("first")
	fmt.Println("  ", result.String())

	// 5. Multiple equivalent ways to get the same value
	fmt.Println("\n5. Different ways to get 'name.last':")
	fmt.Println("   gyaml.Get(yaml, \"name.last\"):", gyaml.Get(yaml, "name.last").String())
	fmt.Println("   gyaml.Get(yaml, \"name\").Get(\"last\"):", gyaml.Get(yaml, "name").Get("last").String())
	fmt.Println("   gyaml.Parse(yaml).Get(\"name\").Get(\"last\"):", gyaml.Parse(yaml).Get("name").Get("last").String())

	// 6. Get string value
	fmt.Println("\n6. Get city:")
	city := gyaml.Get(yaml, "city")
	fmt.Println("  ", city.String())

	// 7. Get entire object as raw YAML
	fmt.Println("\n7. Get name as raw YAML:")
	nameObj := gyaml.Get(yaml, "name")
	fmt.Println("   Raw:", nameObj.Raw)
	fmt.Println("   Type:", nameObj.Type)
}

func main() {
	BasicExample()
}
