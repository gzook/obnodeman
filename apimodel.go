package main

import (
	"encoding/json"
)

type apiSimpleResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r *apiSimpleResponse) ToJSON() *[]byte {
	b, err := json.Marshal(r)
	if err != nil {
		b = []byte(`{"success": false, "error": "failed to convert response to JSON"}`)
	}
	return &b
}
