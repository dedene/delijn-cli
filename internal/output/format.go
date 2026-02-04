package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Mode represents the output mode.
type Mode int

const (
	ModeTable Mode = iota
	ModePlain
	ModeJSON
)

// JSON outputs v as indented JSON to stdout.
func JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}
