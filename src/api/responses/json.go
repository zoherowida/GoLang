package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON displays a json response message with data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintln(w, "%s", err.Error())
	}
}

// ERROR displays a json response message with an error
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Data    string `json:"data"`
			Error   string `json:"error"`
		}{
			Status:  statusCode,
			Message: "error",
			Data:    "",
			Error:   err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
