package image

import "container-inspector/internal/registry"

// LayerInfo combines a layer's size/digest with the Dockerfile command that created it
type LayerInfo struct {
	Index      int
	Digest     string
	Size       int64
	CreatedBy  string
	EmptyLayer bool
}

// Summary holds the full analysis of an image
type Summary struct {
	Image        string
	Tag          string
	Architecture string
	OS           string
	TotalSize    int64
	Layers       []LayerInfo
}

// BuildSummary correlates the manifest layers with the image config history
// History entries with empty_layer=true don't have a corresponding layer in the manifest,
// so we track a separate layer index that only increments for real layers
func BuildSummary(image, tag string, manifest *registry.Manifest, config *registry.ImageConfig) *Summary {
	summary := &Summary{
		Image:        image,
		Tag:          tag,
		Architecture: config.Architecture,
		OS:           config.OS,
	}

	layerIndex := 0

	// Loop through the history entries in the image config and correlate them with the manifest layers
	// For each history entry, we create a LayerInfo struct that includes the created_by command and whether it's an empty layer
	// If the history entry is not an empty layer, we also add the corresponding layer information from the manifest (digest and size) and increment the layer index
	// We also keep a running total of the image size by summing the sizes of the real layers from the manifest
	for _, h := range config.History {
		info := LayerInfo{
			CreatedBy:  h.CreatedBy,
			EmptyLayer: h.EmptyLayer,
		}

		if !h.EmptyLayer && layerIndex < len(manifest.Layers) {
			layer := manifest.Layers[layerIndex]
			info.Index = layerIndex + 1
			info.Digest = layer.Digest
			info.Size = layer.Size
			summary.TotalSize += layer.Size
			layerIndex++
		}

		summary.Layers = append(summary.Layers, info)
	}

	// Return the final summary struct that contains all the correlated information about the image layers and their corresponding Dockerfile commands
	return summary
}
