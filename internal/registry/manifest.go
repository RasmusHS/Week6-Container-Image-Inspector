package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchManifest retrieves the image manifest from the registry using the auth token
func (rc *RegistryClient) FetchManifest(image, tag, token string) (*Manifest, error) {
	body, err := rc.fetchManifestRaw(image, tag, token)
	if err != nil {
		return nil, err
	}

	// Check if this is a manifest index (multi-arch) by peeking at the mediaType
	var probe struct {
		MediaType string `json:"mediaType"`
	}
	json.Unmarshal(body, &probe)

	// If it's a manifest index, we need to find the linux/amd64 entry and fetch that manifest instead
	if probe.MediaType == "application/vnd.oci.image.index.v1+json" ||
		probe.MediaType == "application/vnd.docker.distribution.manifest.list.v2+json" {
		// It's a manifest index — find the amd64/linux entry
		var index ManifestIndex

		// Decode the manifest index JSON into a ManifestIndex struct
		// If there was an error decoding the JSON, return an error indicating that we failed to decode the manifest index
		if err := json.Unmarshal(body, &index); err != nil {
			return nil, fmt.Errorf("failed to decode manifest index: %w", err)
		}

		// Loop through the index entries to find the one that matches linux/amd64 and get its digest
		// If we find a matching entry, we store its digest; if we don't find any matching entry, we return an error indicating that no linux/amd64 manifest was found in the index
		digest := ""
		for _, entry := range index.Manifests {
			if entry.Platform.Architecture == "amd64" && entry.Platform.OS == "linux" {
				digest = entry.Digest
				break
			}
		}
		if digest == "" {
			return nil, fmt.Errorf("no linux/amd64 manifest found in index")
		}

		// Fetch the actual manifest using the digest
		body, err = rc.fetchManifestRaw(image, digest, token)
		if err != nil {
			return nil, err
		}
	}

	// Decode the JSON response body into a Manifest struct and return it
	// If there was an error decoding the JSON, return an error indicating that we failed to decode the manifest
	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	return &manifest, nil
}

// FetchImageConfig retrieves the image configuration blob using the config digest from the manifest
// This contains the history array with Dockerfile commands for each layer
func (rc *RegistryClient) FetchImageConfig(image, token string, configDigest string) (*ImageConfig, error) {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", rc.BaseURL, image, configDigest)

	// Create a new GET request to the blob URL with the Authorization header set to the Bearer token
	// This will fetch the image configuration blob from the registry, which contains the history of Dockerfile commands that created each layer
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the Bearer token to authenticate the request to the registry
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request to the registry to fetch the image configuration blob for the specified image and config digest using the provided auth token
	resp, err := rc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	// If it's not 200, return an error indicating that we failed to fetch the image configuration blob
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image config: status %d", resp.StatusCode)
	}

	// Decode the JSON response body into an ImageConfig struct and return it
	// If there was an error decoding the JSON, return an error indicating that we failed to decode the image configuration blob
	var config ImageConfig
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %w", err)
	}

	// Return the decoded ImageConfig struct, which contains the history of Dockerfile commands that created each layer
	return &config, nil
}

// fetchManifestRaw makes the HTTP request and returns the raw body bytes
func (rc *RegistryClient) fetchManifestRaw(image, reference, token string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.BaseURL, image, reference)

	// Create a new GET request to the manifest URL with the Authorization header set to the Bearer token
	// This will fetch the manifest from the registry, which contains metadata about the image layers and the config blob
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the Bearer token to authenticate the request to the registry
	// Set the Accept header to specify the manifest media types we support (Docker v2, OCI, and multi-arch lists)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
	}, ", "))

	// Make the request to the registry to fetch the manifest for the specified image and reference (tag or digest) using the provided auth token
	// This will return the raw bytes of the manifest JSON, which we can then decode into a Manifest struct
	// If there was an error making the request, return the error; if the status code is not 200 OK, return an error indicating that we failed to fetch the manifest with the given status code
	resp, err := rc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	// If it's not 200, return an error indicating that we failed to fetch the manifest with the given status code
	// If it is 200, read and return the raw bytes of the manifest JSON from the response body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest: status %d", resp.StatusCode)
	}

	// Read and return the raw bytes of the manifest JSON from the response body
	return io.ReadAll(resp.Body)
}
