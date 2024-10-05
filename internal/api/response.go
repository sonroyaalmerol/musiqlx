package api

import (
	"encoding/json"
	"net/http"

	"github.com/sonroyaalmerol/musiqlx/pkg/utils"
)

var SubsonicSuccessResponse = GenericSubsonicResponse{
	Status:        "ok",
	Version:       "1.16.1",
	Type:          "MusiQLx",
	ServerVersion: "0.1.3 (tag)",
	OpenSubsonic:  true,
}

// Centralized success handling function
func HandleSuccess[T any](w http.ResponseWriter, data T) {
	response := map[string]interface{}{
		"subsonic-response": SubsonicSuccessResponse,
	}

	mappedData := utils.StructToMap(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.MergeMaps(response, mappedData))
}
