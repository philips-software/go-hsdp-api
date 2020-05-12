package fhir

import (
	"fmt"
	"sort"
	"strings"
)

// ParseError returns a string representation of an error struct
func ParseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, ParseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, ParseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")
	case float64:
		return fmt.Sprintf("%d", int64(raw))
	case int64:
		return fmt.Sprintf("%d", raw)
	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}
