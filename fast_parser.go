// Copyright 2024 GYAML Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gyaml

import (
	"strconv"
	"strings"
)

// fastGet implements a high-performance direct YAML parser for simple paths
// This is the fast path that avoids yaml.Unmarshal for common cases
func fastGet(yaml, path string) (Result, bool) {
	if len(path) == 0 {
		return Result{Type: YAML, Raw: yaml}, true
	}

	// Parse the path into components
	parts := splitPath(path)
	if len(parts) == 0 {
		return Result{}, false
	}

	// Check if this is a complex path that needs slow path
	if hasComplexFeatures(path) {
		return Result{}, false
	}

	// Check if path ends with getting a collection (no final key)
	// Like "children" or "friends" which should return arrays/objects
	// These need slow path for proper Array() and Map() support
	lastPart := parts[len(parts)-1]
	if lastPart != "#" {
		// For non-terminal paths that retrieve collections, use slow path
		// We'll only use fast path for scalar value retrieval
	}

	// Start parsing from the beginning
	result, ok := fastParsePath(yaml, parts, 0)
	if !ok {
		return Result{}, false
	}

	// If result is a collection (array/object), fall back to slow path
	// for proper support of Array(), Map(), ForEach() etc.
	if result.Type == YAML || result.Type == Null {
		return Result{}, false
	}

	return result, true
}

// hasComplexFeatures checks if the path requires the slow path
func hasComplexFeatures(path string) bool {
	// Check for features that need slow path:
	// - Wildcards: *, ?
	// - Queries: #(...)
	// - Modifiers: @...
	// - Pipes: |
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case '*', '?', '@', '|':
			return true
		case '#':
			if i+1 < len(path) && path[i+1] == '(' {
				return true
			}
		}
	}
	return false
}

// splitPath splits a path by dots, handling escapes
func splitPath(path string) []string {
	if path == "" {
		return nil
	}

	var parts []string
	var current strings.Builder
	escaped := false

	for i := 0; i < len(path); i++ {
		ch := path[i]
		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '.' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// fastParsePath recursively parses YAML following the path
func fastParsePath(yaml string, parts []string, depth int) (Result, bool) {
	if depth >= len(parts) {
		// We've consumed all path parts
		value, ok := extractValue(yaml)
		return value, ok
	}

	currentKey := parts[depth]

	// Check if this is an array index
	if idx, err := strconv.Atoi(currentKey); err == nil {
		return fastParseArrayIndex(yaml, parts, depth, idx)
	}

	// Check if this is a count operation
	if currentKey == "#" {
		// If there are more parts after #, we need to use slow path
		// because # with subsequent parts means "apply path to all elements"
		if depth+1 < len(parts) {
			return Result{}, false
		}
		// Simple count operation
		count := countArrayElements(yaml)
		return Result{
			Type: Number,
			Num:  float64(count),
			Raw:  strconv.Itoa(count),
		}, true
	}

	// It's a key lookup
	return fastParseKey(yaml, parts, depth, currentKey)
}

// fastParseKey finds a key in YAML and continues parsing
func fastParseKey(yaml string, parts []string, depth int, key string) (Result, bool) {
	// Try to find the key in the YAML
	lines := strings.Split(yaml, "\n")

	// Determine base indentation
	baseIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if baseIndent == -1 {
			baseIndent = len(line) - len(strings.TrimLeft(line, " \t"))
		}
		break
	}

	targetIndent := baseIndent
	keyWithColon := key + ":"

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Calculate indentation
		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		// Skip if wrong indentation level
		if indent != targetIndent {
			continue
		}

		// Check if this line contains our key
		if !strings.HasPrefix(trimmed, keyWithColon) {
			continue
		}

		// Found the key! Extract the value part
		valuePart := strings.TrimSpace(trimmed[len(keyWithColon):])

		if depth == len(parts)-1 {
			// This is the final key, extract the value
			if valuePart != "" {
				// Inline value (flow style)
				// Check if it looks like a collection - if so, use slow path
				trimmed := strings.TrimSpace(valuePart)
				if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") ||
					strings.Contains(valuePart, ":") {
					// Likely a collection, fall back to slow path
					return Result{}, false
				}
				return extractValue(valuePart)
			}

			// Block style - value is on next lines with more indentation
			// For block style, if the final value is a collection (array/object),
			// we should fall back to slow path for proper handling
			var blockLines []string
			nextIndent := indent + 2 // YAML standard is 2 spaces

			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]
				nextTrimmed := strings.TrimSpace(nextLine)

				if nextTrimmed == "" {
					continue
				}

				nextLineIndent := len(nextLine) - len(strings.TrimLeft(nextLine, " \t"))

				if nextLineIndent <= indent {
					// Back to same or less indentation, we're done
					break
				}

				if nextLineIndent >= nextIndent {
					blockLines = append(blockLines, nextLine)
				}
			}

			if len(blockLines) > 0 {
				// Check if this is a collection (array or object)
				firstBlockLine := strings.TrimSpace(blockLines[0])
				if strings.HasPrefix(firstBlockLine, "- ") || strings.Contains(firstBlockLine, ": ") {
					// This is a collection, fall back to slow path
					return Result{}, false
				}
				blockValue := strings.Join(blockLines, "\n")
				return extractValue(blockValue)
			}

			return Result{}, false
		}

		// Not the final key, need to recurse
		if valuePart != "" {
			// Inline nested object
			return fastParsePath(valuePart, parts, depth+1)
		}

		// Block style nested object
		var nestedLines []string
		nextIndent := indent + 2

		for j := i + 1; j < len(lines); j++ {
			nextLine := lines[j]
			nextTrimmed := strings.TrimSpace(nextLine)

			if nextTrimmed == "" {
				continue
			}

			nextLineIndent := len(nextLine) - len(strings.TrimLeft(nextLine, " \t"))

			if nextLineIndent <= indent {
				break
			}

			if nextLineIndent >= nextIndent {
				// Adjust indentation to make it relative
				adjusted := strings.Repeat(" ", nextLineIndent-nextIndent) + strings.TrimLeft(nextLine, " \t")
				nestedLines = append(nestedLines, adjusted)
			}
		}

		if len(nestedLines) > 0 {
			nestedYAML := strings.Join(nestedLines, "\n")
			return fastParsePath(nestedYAML, parts, depth+1)
		}

		return Result{}, false
	}

	return Result{}, false
}

// fastParseArrayIndex handles array index access
func fastParseArrayIndex(yaml string, parts []string, depth int, index int) (Result, bool) {
	elements := parseArrayElements(yaml)

	if index < 0 || index >= len(elements) {
		return Result{}, false
	}

	element := elements[index]

	if depth == len(parts)-1 {
		// Final key, extract value
		return extractValue(element)
	}

	// Continue parsing
	return fastParsePath(element, parts, depth+1)
}

// parseArrayElements extracts array elements from YAML
func parseArrayElements(yaml string) []string {
	var elements []string
	lines := strings.Split(yaml, "\n")

	var currentElement strings.Builder
	inElement := false
	baseIndent := -1
	elementIndent := -1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		// Determine base indentation from first non-empty line
		if baseIndent == -1 {
			baseIndent = indent
		}

		// Check if this is a top-level array item (starts with -)
		if strings.HasPrefix(trimmed, "- ") && indent == baseIndent {
			// Save previous element if any
			if inElement && currentElement.Len() > 0 {
				elements = append(elements, currentElement.String())
				currentElement.Reset()
			}

			// Set element indent (the indent of items within this array element)
			elementIndent = indent

			// Extract the value after the dash
			value := strings.TrimSpace(trimmed[2:])
			currentElement.WriteString(value)
			inElement = true
		} else if inElement && indent > elementIndent {
			// This is a continuation of the current element
			// (nested content or multi-line values)
			if currentElement.Len() > 0 {
				currentElement.WriteString("\n")
			}
			currentElement.WriteString(line)
		} else if inElement && indent <= elementIndent {
			// We've reached something at the same or lower indentation
			// which means we might be done with this element
			// But only if it's not another array element at base level
			if strings.HasPrefix(trimmed, "- ") && indent == baseIndent {
				// This is the next array element, we'll handle it in the next iteration
				// Save current and reset
				elements = append(elements, currentElement.String())
				currentElement.Reset()

				elementIndent = indent
				value := strings.TrimSpace(trimmed[2:])
				currentElement.WriteString(value)
				inElement = true
			}
		}
	}

	// Don't forget the last element
	if currentElement.Len() > 0 {
		elements = append(elements, currentElement.String())
	}

	return elements
}

// countArrayElements counts elements in a YAML array
func countArrayElements(yaml string) int {
	return len(parseArrayElements(yaml))
}

// extractValue converts a YAML value string to a Result
func extractValue(value string) (Result, bool) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Result{}, false
	}

	// Handle quoted strings
	if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		unquoted := value[1 : len(value)-1]
		return Result{
			Type: String,
			Str:  unquoted,
			Raw:  value,
		}, true
	}

	// Handle booleans
	lower := strings.ToLower(value)
	if lower == "true" || lower == "yes" || lower == "on" {
		return Result{
			Type: True,
			Raw:  value,
		}, true
	}
	if lower == "false" || lower == "no" || lower == "off" {
		return Result{
			Type: False,
			Raw:  value,
		}, true
	}

	// Handle null
	if lower == "null" || lower == "~" || value == "" {
		return Result{
			Type: Null,
			Raw:  value,
		}, true
	}

	// Try to parse as number
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return Result{
			Type: Number,
			Num:  num,
			Raw:  value,
		}, true
	}

	// Default to string
	return Result{
		Type: String,
		Str:  value,
		Raw:  value,
	}, true
}
