package output

import (
	"encoding/json"
	"strings"
)

// FieldRow represents one row of extracted field data.
type FieldRow struct {
	Fields map[string]string `json:"fields"`
}

// FieldsResult holds parsed tshark field extraction output.
type FieldsResult struct {
	Rows      []FieldRow `json:"rows"`
	Total     int        `json:"total"`
	FieldNames []string  `json:"field_names"`
}

// ParseFieldOutput parses tab-separated tshark -T fields output into structured JSON.
// fieldNames are the -e field names in order.
func ParseFieldOutput(data []byte, fieldNames []string) ([]byte, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return json.Marshal(FieldsResult{Rows: []FieldRow{}, FieldNames: fieldNames})
	}

	rows := make([]FieldRow, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		values := strings.Split(line, "\t")
		row := FieldRow{Fields: make(map[string]string, len(fieldNames))}
		for i, name := range fieldNames {
			if i < len(values) {
				row.Fields[name] = values[i]
			} else {
				row.Fields[name] = ""
			}
		}
		rows = append(rows, row)
	}

	result := FieldsResult{
		Rows:       rows,
		Total:      len(rows),
		FieldNames: fieldNames,
	}

	return json.MarshalIndent(result, "", "  ")
}
