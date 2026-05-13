package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RegistryClient struct {
	Client  *http.Client
	BaseURL string
}

// Constructor for RegistryClient that returns a new instance with the provided base URL and a default HTTP client
func NewClient(baseURL string) *RegistryClient {
	return &RegistryClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

// Method that takes an image and tag, does the auth dance, and returns the token
func (rc *RegistryClient) GetAuthToken(image, tag string) (string, error) {
	// Calls GET to /v2/{image}/manifests/{tag} to get the status code and WWW-Authenticate header
	// Prints the status code and WWW-Authenticate header for debugging purposes
	req, err := http.NewRequest("GET", rc.BaseURL+"/v2/"+image+"/manifests/"+tag, nil)
	if err != nil {
		return "", err
	}

	// Make the request to the registry to get the status code and WWW-Authenticate header for the specified image and tag
	resp, err := rc.Client.Do(req)
	if err != nil {
		return "", err
	}

	// Close the response body when we're done with it to prevent resource leaks
	defer resp.Body.Close()

	// Check if the status code is 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode) // No auth needed, return empty token
	}

	// Get the WWW-Authenticate header from the response to extract the realm, service, and scope for the token request
	wwwAuth := resp.Header.Get("WWW-Authenticate")

	// Parse the WWW-Authenticate header to extract the realm, service, and scope
	realm, service, scope, err := parseWWWAuthenticate(wwwAuth)
	if err != nil {
		return "", err // Return an error if there was an issue parsing the WWW-Authenticate header
	}
	// Make a GET request to the realm URL with the service and scope as query parameters to get the token
	tokenReq, err := http.NewRequest("GET", realm, nil)
	if err != nil {
		return "", err // Return an error if there was an issue creating the token request
	}
	// Add the service and scope as query parameters to the token request
	q := tokenReq.URL.Query() // Add the service and scope as query parameters to the token request
	q.Add("service", service)
	q.Add("scope", scope)

	// Encode the query parameters and set them in the token request URL
	tokenReq.URL.RawQuery = q.Encode()
	tokenResp, err := rc.Client.Do(tokenReq)
	if err != nil {
		return "", err // Return an error if there was an issue making the token request
	}

	// Close the token response body when we're done with it to prevent resource leaks
	defer tokenResp.Body.Close()

	// Parse the token response to extract the token and return it
	var tokenResponse TokenResponse
	err = json.NewDecoder(tokenResp.Body).Decode(&tokenResponse) // Decode the JSON response body into the TokenResponse struct
	if err != nil {
		return "", err // Return an error if there was an issue decoding the JSON response
	}
	return tokenResponse.Token, nil // Return the token from the token response
}

// Helper function to parse the WWW-Authenticate header and extract the realm, service, and scope
func parseWWWAuthenticate(header string) (realm, service, scope string, err error) {
	// Example header: Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:library/alpine:pull"

	// Split the header by spaces to separate the scheme and parameters
	parts := strings.SplitN(header, " ", 2)
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid WWW-Authenticate header: %s", header) // Invalid header format
	}

	// The parameters are in the second part, split by commas
	params := strings.Split(parts[1], ",")
	for _, param := range params {
		// Split each parameter by the equals sign to separate the key and value
		kv := strings.SplitN(param, "=", 2)
		if len(kv) < 2 {
			continue // Invalid parameter format, skip
		}
		key := strings.Trim(kv[0], `"`)   // Trim any surrounding quotes from the key
		value := strings.Trim(kv[1], `"`) // Trim any surrounding quotes from the value
		switch key {
		case "realm": // If the key is "realm", set the realm variable to the value
			realm = value
		case "service": // If the key is "service", set the service variable to the value
			service = value
		case "scope": // If the key is "scope", set the scope variable to the value
			scope = value
		}
	}
	return realm, service, scope, nil // Return the extracted realm, service, and scope values
}
