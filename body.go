package grove

import (
	"encoding/json"
	"net/http"
)

// Helper function that will parse a JSON body from the `*http.Request`.
// If the decoder fails to parse the body into the provided type `T` it will return the
// decoder error.
// This uses the standard `encoding/json` json decoder.
func ParseJsonBodyFromRequest[T any](request *http.Request) (T, error) {
	var body T
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

// Helper function to write something to a JSON response.
// It sets the content-type header and writes the provided body `T` to the response body.
// It uses the standard `encoding/json` json encoder and returns the error result it returns.
func WriteJsonBodyToResponse[T any](response http.ResponseWriter, body T) error {
	response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(response).Encode(body)
}

// Helper function to return an error response.
// It sets the content-type header to `application/json` and returns a json response with the provided
// message. The response shape is `{"error": message}`.
// It also sets the status code to the provided `statusCode`.
func WriteErrorToResponse(response http.ResponseWriter, statusCode int, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	_ = json.NewEncoder(response).Encode(map[string]string{"error": message})
}
