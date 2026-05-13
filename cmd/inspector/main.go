package main

import (
	"container-inspector/internal/image"
	"container-inspector/internal/output"
	"container-inspector/internal/registry"

	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: inspector <image> <tag> [--inspect-layer N]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  inspector library/alpine latest")
		fmt.Println("  inspector library/nginx stable")
		fmt.Println("  inspector library/alpine latest --inspect-layer 1")
		os.Exit(1)
	}

	imageName := os.Args[1]
	tag := os.Args[2]

	// Check for optional --inspect-layer flag
	inspectLayer := -1
	if len(os.Args) == 5 && os.Args[3] == "--inspect-layer" {
		n, err := strconv.Atoi(os.Args[4])
		if err != nil || n < 1 {
			fmt.Println("Error: --inspect-layer must be a positive number")
			os.Exit(1)
		}
		inspectLayer = n
	}

	// Create registry client and authenticate
	client := registry.NewClient("https://registry-1.docker.io")

	fmt.Printf("Fetching manifest for %s:%s...\n\n", imageName, tag)

	// Get the auth token for the image and tag
	token, err := client.GetAuthToken(imageName, tag)
	if err != nil {
		fmt.Printf("Error getting auth token: %v\n", err)
		os.Exit(1)
	}

	// Fetch the image manifest
	manifest, err := client.FetchManifest(imageName, tag, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Manifest error: %v\n", err)
		os.Exit(1)
	}

	// Fetch the image config (contains Dockerfile history)
	config, err := client.FetchImageConfig(imageName, token, manifest.Config.Digest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	// Build and display the summary
	summary := image.BuildSummary(imageName, tag, manifest, config)
	output.PrintSummary(summary)

	// Optional: inspect a specific layer's contents
	if inspectLayer > 0 {
		// Find the layer by index
		for _, layer := range summary.Layers {
			if layer.Index == inspectLayer {
				fmt.Printf("Downloading and inspecting layer %d...\n", inspectLayer)
				entries, err := image.InspectLayer(
					"https://registry-1.docker.io",
					imageName,
					token,
					layer.Digest,
				)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Layer inspect error: %v\n", err)
					os.Exit(1)
				}
				output.PrintLayerContents(inspectLayer, layer.Digest, entries)
				break
			}
		}
	}
}
