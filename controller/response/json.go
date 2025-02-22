package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
	_, err = w.Write(data)
	return
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJson(w, statusCode, map[string]string{
		"error": message,
	})
}
