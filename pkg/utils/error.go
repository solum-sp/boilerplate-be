package utils

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
)

type CustomError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var _ error = (*CustomError)(nil)

func (e *CustomError) Error() string {
	return e.Message
}

var errorManager = make(map[string]*CustomError)

func LoadErrorMessages(fs embed.FS, fileName string) error {
	log.Printf("Loading error messages from %s", fileName)
	data, err := fs.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read errors.json: %w", err)
	}

	var errorsMap map[string]string
	err = json.Unmarshal(data, &errorsMap)
	if err != nil {
		return fmt.Errorf("failed to parse errors.json: %w", err)
	}

	for code, msg := range errorsMap {
		customErr := errorManager[code]
		customErr.Message = msg
	}
	for code, err := range errorManager {
		if err.Message == "" {
			return fmt.Errorf("error %s has no message", code)
		}

	}

	return nil
}

func NewCustomError(code string) *CustomError {
	err := CustomError{
		Code:    code,
		Message: "",
	}
	errorManager[code] = &err
	return &err
}