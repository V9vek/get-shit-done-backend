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
	// TODO: create middleware to set `content type: application/json` on every outgoing response
	write.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(write)
	err := encoder.Encode(response)
	PanicIfError(err)
}
