package main

import (
	"errors"
	"fmt"

	"github.com/m4l1c1ou5/gyaml"
)

// ValidationExample demonstrates YAML validation
func ValidationExample() {
	fmt.Println("=== Validation Examples ===\n")

	// 1. Valid YAML
	const validYaml = `
name: John Doe
age: 30
active: true
tags:
  - golang
  - yaml
`
	fmt.Println("1. Valid YAML:")
	fmt.Println(validYaml)
	fmt.Println("   Is valid:", gyaml.Valid(validYaml))

	// 2. Invalid YAML - unclosed bracket
	const invalidYaml1 = `
name: John Doe
tags:
  - golang
  - yaml
  - [unclosed
`
	fmt.Println("\n2. Invalid YAML (unclosed bracket):")
	fmt.Println(invalidYaml1)
	fmt.Println("   Is valid:", gyaml.Valid(invalidYaml1))

	// 3. Invalid YAML - bad indentation
	const invalidYaml2 = `
name: John Doe
  age: 30
active: true
`
	fmt.Println("\n3. Invalid YAML (bad indentation):")
	fmt.Println(invalidYaml2)
	fmt.Println("   Is valid:", gyaml.Valid(invalidYaml2))

	// 4. Invalid YAML - incorrect map
	const invalidYaml3 = `
name: John
  - item1
  - item2
`
	fmt.Println("\n4. Invalid YAML (mixed map/array):")
	fmt.Println(invalidYaml3)
	fmt.Println("   Is valid:", gyaml.Valid(invalidYaml3))

	// 5. Empty YAML
	const emptyYaml = ``
	fmt.Println("\n5. Empty YAML:")
	fmt.Println("   Is valid:", gyaml.Valid(emptyYaml))

	// 6. YAML with only whitespace
	const whitespaceYaml = `
	
	
`
	fmt.Println("\n6. Whitespace only YAML:")
	fmt.Println("   Is valid:", gyaml.Valid(whitespaceYaml))

	// 7. Using Valid() before Get() - best practice
	fmt.Println("\n7. Best practice - validate before use:")
	yamlInput := `
user:
  name: Alice
  email: alice@example.com
`
	if err := validateAndProcess(yamlInput); err != nil {
		fmt.Println("   Error:", err)
	} else {
		fmt.Println("   Processing completed successfully")
	}

	// 8. Using Valid() with untrusted input
	fmt.Println("\n8. Validate untrusted input:")
	untrustedInput := `
data:
  key: value
  nested:
    item: test
`
	if gyaml.Valid(untrustedInput) {
		result := gyaml.Get(untrustedInput, "data.nested.item")
		fmt.Println("   Retrieved value:", result.String())
	} else {
		fmt.Println("   Invalid YAML from untrusted source")
	}

	// 9. Complex valid YAML
	const complexYaml = `
company:
  name: TechCorp
  employees:
    - name: John
      department: Engineering
      skills:
        - Go
        - Python
    - name: Jane
      department: Marketing
      skills:
        - SEO
        - Content
  locations:
    - city: NYC
      country: USA
    - city: London
      country: UK
`
	fmt.Println("\n9. Complex valid YAML:")
	fmt.Println("   Is valid:", gyaml.Valid(complexYaml))

	// 10. Using @valid modifier
	fmt.Println("\n10. Using @valid modifier:")
	validResult := gyaml.Get(validYaml, "@valid")
	fmt.Println("   Result:", validResult.Raw)

	// 11. Multi-document YAML
	const multiDoc = `
---
doc1:
  name: First
---
doc2:
  name: Second
`
	fmt.Println("\n11. Multi-document YAML:")
	fmt.Println("   Is valid:", gyaml.Valid(multiDoc))

	// 12. YAML with special characters
	const specialChars = `
text: "Hello \"World\""
path: "C:\\Users\\John"
unicode: "こんにちは"
`
	fmt.Println("\n12. YAML with special characters:")
	fmt.Println("   Is valid:", gyaml.Valid(specialChars))
}

// validateAndProcess demonstrates proper validation pattern
func validateAndProcess(yaml string) error {
	if !gyaml.Valid(yaml) {
		return errors.New("invalid yaml")
	}

	// Safe to process
	name := gyaml.Get(yaml, "user.name")
	email := gyaml.Get(yaml, "user.email")

	fmt.Printf("   User: %s (%s)\n", name.String(), email.String())
	return nil
}

func main() {
	ValidationExample()
}
