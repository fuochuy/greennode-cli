package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmespath/go-jmespath"
)

// Format formats and outputs the response data.
func Format(data interface{}, outputFormat, query string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	// Apply JMESPath query if specified
	if query != "" {
		result, err := jmespath.Search(query, data)
		if err != nil {
			return fmt.Errorf("JMESPath query error: %w", err)
		}
		data = result
	}

	if data == nil {
		return nil
	}

	switch outputFormat {
	case "json":
		formatJSON(data, w)
	case "text":
		formatText(data, w)
	case "table":
		formatTable(data, w)
	default:
		formatJSON(data, w)
	}
	return nil
}

func formatJSON(data interface{}, w io.Writer) {
	if isEmptyMap(data) {
		return
	}
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Fprintf(w, "%v\n", data)
		return
	}
	fmt.Fprintf(w, "%s\n", string(out))
}

func formatText(data interface{}, w io.Writer) {
	if data == nil || isEmptyMap(data) {
		return
	}

	switch v := data.(type) {
	case map[string]interface{}:
		for _, value := range v {
			if items, ok := value.([]interface{}); ok {
				for _, item := range items {
					if m, ok := item.(map[string]interface{}); ok {
						vals := mapValues(m)
						fmt.Fprintln(w, strings.Join(vals, "\t"))
					} else {
						fmt.Fprintln(w, fmt.Sprint(item))
					}
				}
				return
			}
		}
		vals := mapValues(v)
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				vals := mapValues(m)
				fmt.Fprintln(w, strings.Join(vals, "\t"))
			} else {
				fmt.Fprintln(w, fmt.Sprint(item))
			}
		}
	default:
		fmt.Fprintln(w, fmt.Sprint(data))
	}
}

func formatTable(data interface{}, w io.Writer) {
	if data == nil || isEmptyMap(data) {
		return
	}

	rows := extractRows(data)
	if len(rows) == 0 {
		return
	}

	// Check if rows are maps
	if firstMap, ok := rows[0].(map[string]interface{}); ok {
		headers := mapKeys(firstMap)
		colWidths := make([]int, len(headers))
		for i, h := range headers {
			colWidths[i] = len(h)
		}

		strRows := make([][]string, len(rows))
		for i, row := range rows {
			m, _ := row.(map[string]interface{})
			strRows[i] = make([]string, len(headers))
			for j, h := range headers {
				val := fmt.Sprint(m[h])
				strRows[i][j] = val
				if len(val) > colWidths[j] {
					colWidths[j] = len(val)
				}
			}
		}

		// Print header
		headerParts := make([]string, len(headers))
		sepParts := make([]string, len(headers))
		for i, h := range headers {
			headerParts[i] = padRight(h, colWidths[i])
			sepParts[i] = strings.Repeat("-", colWidths[i])
		}
		fmt.Fprintln(w, strings.Join(headerParts, " | "))
		fmt.Fprintln(w, strings.Join(sepParts, "-+-"))

		// Print rows
		for _, row := range strRows {
			parts := make([]string, len(row))
			for i, val := range row {
				parts[i] = padRight(val, colWidths[i])
			}
			fmt.Fprintln(w, strings.Join(parts, " | "))
		}
	}
}

func extractRows(data interface{}) []interface{} {
	switch v := data.(type) {
	case []interface{}:
		return v
	case map[string]interface{}:
		for _, value := range v {
			if items, ok := value.([]interface{}); ok {
				return items
			}
		}
		return []interface{}{v}
	default:
		return []interface{}{v}
	}
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func mapValues(m map[string]interface{}) []string {
	vals := make([]string, 0, len(m))
	for _, v := range m {
		vals = append(vals, fmt.Sprint(v))
	}
	return vals
}

func isEmptyMap(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	return ok && len(m) == 0
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
