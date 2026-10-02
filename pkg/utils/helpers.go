package utils

import (
	"encoding/json"
	"fmt"
)

func ParseJSON(jsonStr string, out interface{}) error {
	if err := json.Unmarshal([]byte(jsonStr), out); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func DataParser1(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal src: %w", err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("failed to unmarshal to dst: %w", err)
	}
	return nil
}
