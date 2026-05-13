package registry

// TokenResponse represents the JSON response from the auth token endpoint
type TokenResponse struct {
	Token string `json:"token"`
}

// Descriptor is a common structure used for both config and layer references in a manifest
type Descriptor struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// Manifest represents a Docker/OCI image manifest (schema version 2)
type Manifest struct {
	SchemaVersion int          `json:"schemaVersion"`
	MediaType     string       `json:"mediaType"`
	Config        Descriptor   `json:"config"`
	Layers        []Descriptor `json:"layers"`
}

// ImageConfig represents the image configuration blob
// Contains the history of Dockerfile commands that created each layer
type ImageConfig struct {
	Architecture string         `json:"architecture"`
	OS           string         `json:"os"`
	History      []HistoryEntry `json:"history"`
}

// HistoryEntry represents a single step in the image build history
type HistoryEntry struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	EmptyLayer bool   `json:"empty_layer,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// ManifestIndex represents a multi-arch manifest list / OCI image index
type ManifestIndex struct {
	SchemaVersion int          `json:"schemaVersion"`
	MediaType     string       `json:"mediaType"`
	Manifests     []IndexEntry `json:"manifests"`
}

// IndexEntry is a single platform entry in a manifest index
type IndexEntry struct {
	MediaType string   `json:"mediaType"`
	Size      int64    `json:"size"`
	Digest    string   `json:"digest"`
	Platform  Platform `json:"platform"`
}

// Platform describes the OS and architecture of a manifest entry
type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}
