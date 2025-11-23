package main

import (
	"fmt"
	"strings"

	"github.com/m4l1c1ou5/gyaml"
)

// CustomModifiersExample demonstrates creating and using custom modifiers
func CustomModifiersExample() {
	const yaml = `
message: hello world
items:
  - apple
  - banana
  - cherry
name:
  first: john
  last: doe
numbers:
  - 1
  - 2
  - 3
`

	fmt.Println("=== Custom Modifiers Examples ===\n")

	// 1. Custom modifier: case (upper/lower)
	gyaml.AddModifier("case", func(yaml, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(yaml)
		}
		if arg == "lower" {
			return strings.ToLower(yaml)
		}
		return yaml
	})

	fmt.Println("1. Custom 'case:upper' modifier (message|@case:upper):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@case:upper").String())

	fmt.Println("\n2. Custom 'case:lower' modifier (message|@case:lower):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@case:lower").String())

	// 2. Custom modifier: repeat
	gyaml.AddModifier("repeat", func(yaml, arg string) string {
		if arg == "" {
			arg = "2"
		}
		count := 2
		fmt.Sscanf(arg, "%d", &count)
		result := ""
		for i := 0; i < count; i++ {
			if i > 0 {
				result += " "
			}
			result += strings.TrimSpace(yaml)
		}
		return result
	})

	fmt.Println("\n3. Custom 'repeat' modifier (items.0|@repeat:3):")
	fmt.Println("  ", gyaml.Get(yaml, "items.0|@repeat:3").String())

	// 3. Custom modifier: prefix
	gyaml.AddModifier("prefix", func(yaml, arg string) string {
		return arg + yaml
	})

	fmt.Println("\n4. Custom 'prefix' modifier (message|@prefix:greeting: ):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@prefix:greeting: ").String())

	// 4. Custom modifier: suffix
	gyaml.AddModifier("suffix", func(yaml, arg string) string {
		return yaml + arg
	})

	fmt.Println("\n5. Custom 'suffix' modifier (message|@suffix:!):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@suffix:!").String())

	// 5. Custom modifier: title case
	gyaml.AddModifier("title", func(yaml, arg string) string {
		return strings.Title(yaml)
	})

	fmt.Println("\n6. Custom 'title' modifier (message|@title):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@title").String())

	// 6. Chaining custom modifiers
	fmt.Println("\n7. Chain custom modifiers (message|@case:upper|@suffix:!!!):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@case:upper|@suffix:!!!").String())

	// 7. Custom modifier with array
	gyaml.AddModifier("count", func(yaml, arg string) string {
		result := gyaml.Parse(yaml)
		if result.IsArray() {
			return fmt.Sprintf("%d", len(result.Array()))
		}
		return "0"
	})

	fmt.Println("\n8. Custom 'count' modifier on array (items|@count):")
	fmt.Println("  ", gyaml.Get(yaml, "items|@count").String())

	// 8. Custom modifier: reverse string
	gyaml.AddModifier("reversestr", func(yaml, arg string) string {
		runes := []rune(yaml)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})

	fmt.Println("\n9. Custom 'reversestr' modifier (message|@reversestr):")
	fmt.Println("  ", gyaml.Get(yaml, "message|@reversestr").String())

	// 9. Using custom modifier with path continuation
	fmt.Println("\n10. Custom modifier with path (items|@case:upper|0):")
	fmt.Println("  ", gyaml.Get(yaml, "items|@case:upper|0").String())
}

func main() {
	CustomModifiersExample()
}
