package main

import (
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// QueriesExample demonstrates query syntax and operations
func QueriesExample() {
	const yaml = `
friends:
  - first: Dale
    last: Murphy
    age: 44
    nets:
      - ig
      - fb
      - tw
  - first: Roger
    last: Craig
    age: 68
    nets:
      - fb
      - tw
  - first: Jane
    last: Murphy
    age: 47
    nets:
      - ig
      - tw
  - first: David
    last: Smith
    age: 32
    nets:
      - tw
products:
  - name: Laptop
    price: 999.99
    inStock: true
  - name: Mouse
    price: 29.99
    inStock: false
  - name: Keyboard
    price: 79.99
    inStock: true
`

	fmt.Println("=== Query Examples ===\n")

	// 1. Equality query - first match
	fmt.Println("1. First match equality (friends.#(last==\"Murphy\").first):")
	fmt.Println("  ", gyaml.Get(yaml, `friends.#(last=="Murphy").first`).String())

	// 2. Equality query - all matches
	fmt.Println("\n2. All matches equality (friends.#(last==\"Murphy\")#.first):")
	results := gyaml.Get(yaml, `friends.#(last=="Murphy")#.first`)
	for _, r := range results.Array() {
		fmt.Println("  ", r.String())
	}

	// 3. Greater than query
	fmt.Println("\n3. Greater than (friends.#(age>45)#.last):")
	older := gyaml.Get(yaml, "friends.#(age>45)#.last")
	for _, r := range older.Array() {
		fmt.Println("  ", r.String())
	}

	// 4. Less than query
	fmt.Println("\n4. Less than (friends.#(age<40)#.first):")
	younger := gyaml.Get(yaml, "friends.#(age<40)#.first")
	for _, r := range younger.Array() {
		fmt.Println("  ", r.String())
	}

	// 5. Greater than or equal query
	fmt.Println("\n5. Greater than or equal (friends.#(age>=47)#.first):")
	ageGte := gyaml.Get(yaml, "friends.#(age>=47)#.first")
	for _, r := range ageGte.Array() {
		fmt.Println("  ", r.String())
	}

	// 6. Less than or equal query
	fmt.Println("\n6. Less than or equal (friends.#(age<=44)#.first):")
	ageLte := gyaml.Get(yaml, "friends.#(age<=44)#.first")
	for _, r := range ageLte.Array() {
		fmt.Println("  ", r.String())
	}

	// 7. Not equal query
	fmt.Println("\n7. Not equal (friends.#(last!=\"Murphy\")#.first):")
	notMurphy := gyaml.Get(yaml, `friends.#(last!="Murphy")#.first`)
	for _, r := range notMurphy.Array() {
		fmt.Println("  ", r.String())
	}

	// 8. Pattern matching with % (like)
	fmt.Println("\n8. Pattern match like (friends.#(first%\"D*\").last):")
	fmt.Println("  ", gyaml.Get(yaml, `friends.#(first%"D*").last`).String())

	// 9. Pattern matching with !% (not like)
	fmt.Println("\n9. Pattern match not like (friends.#(first!%\"D*\").last):")
	fmt.Println("  ", gyaml.Get(yaml, `friends.#(first!%"D*").last`).String())

	// 10. Nested array query
	fmt.Println("\n10. Nested array query (friends.#(nets.#(==\"fb\"))#.first):")
	hasFb := gyaml.Get(yaml, `friends.#(nets.#(=="fb"))#.first`)
	for _, r := range hasFb.Array() {
		fmt.Println("  ", r.String())
	}

	// 11. Boolean query
	fmt.Println("\n11. Boolean query (products.#(inStock==true)#.name):")
	inStock := gyaml.Get(yaml, "products.#(inStock==true)#.name")
	for _, r := range inStock.Array() {
		fmt.Println("  ", r.String())
	}

	// 12. Float comparison query
	fmt.Println("\n12. Float comparison (products.#(price<50)#.name):")
	affordable := gyaml.Get(yaml, "products.#(price<50)#.name")
	for _, r := range affordable.Array() {
		fmt.Println("  ", r.String())
	}

	// 13. Query with multiple conditions (using nested queries)
	fmt.Println("\n13. Get first person under 40:")
	fmt.Println("  ", gyaml.Get(yaml, `friends.#(age<40).first`).String())
}

func main() {
	QueriesExample()
}
