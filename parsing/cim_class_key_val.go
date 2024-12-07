package parsing

import (
	"fmt"
	"regexp"
	"strings"
)

// CimClassKeyVal represents a map of key-value pairs.
// This is used to represent some CimClass fields that are returned
// as a single string with key-value pairs.
type CimClassKeyVal map[string]string

// UnmarshalJSON unmarshals a JSON object into a map of strings.
// The expected format is a string with key-value pairs separated by spaces, where the key is
// separated from the value by an equals sign.
func (kv *CimClassKeyVal) UnmarshalJSON(b []byte) error {
	// Convert the input bytes to a string
	raw := string(b)

	// Remove surrounding quotes
	if strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
		raw = raw[1 : len(raw)-1]
	}

	// Unescape the JSON string
	raw = strings.ReplaceAll(raw, `\"`, `"`)

	// Initialize the result map.
	result := make(map[string]string)

	// Regular expression to match key-value pairs with quoted and unquoted values.
	pairRegex, err := regexp.Compile(`(\S+)\s*=\s*("(.*?)"|'(.*?)'|(\S+))`)
	if err != nil {
		return fmt.Errorf("parsing.CimClassKeyVal.UnmarshalJSON: %s", err)
	}

	// Find all key-value pairs in the input string.
	matches := pairRegex.FindAllStringSubmatch(raw, -1)

	for _, match := range matches {
		key := match[1]

		// Use the longest group that matches for the value
		value := match[3] // Match from the first quoted group

		if value == "" {
			value = match[4] // Match from the second quoted group
		}
		if value == "" {
			value = match[5] // Match unquoted value
		}

		// Write the key-value pair to the result.
		result[key] = value
	}

	*kv = result

	return nil
}