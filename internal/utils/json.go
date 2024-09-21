package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveJSON(filename string, data interface{}) error {
	// os.Create will overwrite the file if it already exists
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create or overwrite file %s: %v", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding data to %s: %w", filename, err)
	}

	fmt.Printf("File saved successfully: %s\n", filename)
	return nil
}
