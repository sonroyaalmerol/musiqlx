package api

import (
	"encoding/json"
	"net/http"
)

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type SubsonicError struct {
	GenericSubsonicResponse
	Error ErrorDetails `json:"error"`
}

var SubsonicErrorResponse = GenericSubsonicResponse{
	Status:        "failed",
	Version:       "1.16.1",
	Type:          "MusiQLx",
	ServerVersion: "0.1.3 (tag)",
	OpenSubsonic:  true,
}

func HandleError(w http.ResponseWriter, code int, message string) {
	details := ErrorDetails{
		Code:    code,
		Message: message,
	}

	errorResponse := SubsonicError{
		GenericSubsonicResponse: SubsonicErrorResponse,
		Error:                   details,
	}

	response := map[string]interface{}{
		"subsonic-response": errorResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// Predefined error responses for common use cases
func MissingParameterError(w http.ResponseWriter) {
	HandleError(w, 10, "Required parameter is missing.")
}

func WrongCredentialsError(w http.ResponseWriter) {
	HandleError(w, 40, "Wrong username or password.")
}

func UnauthorizedError(w http.ResponseWriter) {
	HandleError(w, 50, "User is not authorized for the given operation.")
}

func NotFoundError(w http.ResponseWriter) {
	HandleError(w, 70, "The requested data was not found.")
}
