package utils

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, status int, body interface{}) {
	if body != nil {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	if body != nil {
		jsn, err := json.Marshal(body)
		if err != nil {
			return
		}
		_, _ = w.Write(jsn)
	}
}
