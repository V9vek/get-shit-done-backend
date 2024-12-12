package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadFromRequestBody(r *http.Request, result interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return fmt.Errorf("not able to parse data from request body: %w", err)
	}
	return nil
}

func WriteResponseBody(write http.ResponseWriter, response interface{}) {
	write.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(write).Encode(response)
	PanicIfError(err)
}
