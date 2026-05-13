package image

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// LayerEntry represents a single file/directory inside a layer
type LayerEntry struct {
	Name string
	Size int64
	Type string // "file", "dir", "symlink"
}

// InspectLayer downloads a layer blob from the registry and lists its contents
// This is the gzip/tar inspection part — each layer is a .tar.gz archive
func InspectLayer(registryURL, image, token, digest string) ([]LayerEntry, error) {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", registryURL, image, digest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch layer: status %d", resp.StatusCode)
	}

	// Layers are gzip-compressed tar archives
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	var entries []LayerEntry

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading tar: %w", err)
		}

		entry := LayerEntry{
			Name: header.Name,
			Size: header.Size,
		}

		switch header.Typeflag {
		case tar.TypeDir:
			entry.Type = "dir"
		case tar.TypeSymlink, tar.TypeLink:
			entry.Type = "symlink"
		default:
			entry.Type = "file"
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
