package registry

import (
	"net/http"
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
	resp, err := rc.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the status code is 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		return "", nil // No auth needed, return empty token
	}
	wwwAuth := resp.Header.Get("WWW-Authenticate")
	println("Status Code:", resp.StatusCode)
	println("WWW-Authenticate Header:", wwwAuth)

	// Dummy return for now, we will implement the actual token retrieval logic later
	return "dummy-token", nil
}
