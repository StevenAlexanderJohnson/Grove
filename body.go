package grove

import (
	"encoding/json"
	"net/http"
)

func ParseJsonBodyFromRequest[T any](request *http.Request) (T, error) {
	var body T
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

func WriteJsonBodyToResponse[T any](response http.ResponseWriter, body T) error {
	response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(response).Encode(body)
}

func WriteErrorToResponse(response http.ResponseWriter, statusCode int, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	_ = json.NewEncoder(response).Encode(map[string]string{"error": message})
}
