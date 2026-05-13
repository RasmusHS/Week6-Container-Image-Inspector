package main

import (
	"container-inspector/internal/registry"
	"fmt"
	"os"
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

	image := os.Args[1]
	tag := os.Args[2]

	fmt.Printf("Image: %s, Tag: %s\n", image, tag)

	// Call the registry package to get the client
	client := registry.NewClient("https://registry-1.docker.io")

	// Get the auth token for the image and tag
	token, err := client.GetAuthToken(image, tag)
	if err != nil {
		fmt.Printf("Error getting auth token: %v\n", err)
		os.Exit(1)
	}

	// Fetch the image manifest
	manifest, err := client.FetchManifest(image, tag, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Manifest error: %v\n", err)
		os.Exit(1)
	}

	// Debug print the config digest from the manifest for verification
	fmt.Printf("Config digest: %s\n", manifest.Config.Digest)

	// Fetch the image config (contains Dockerfile history)
	config, err := client.FetchImageConfig(image, token, manifest.Config.Digest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	// Debug print the auth token, manifest and config for verification
	fmt.Printf("Auth Token: %s\n", token)
	fmt.Printf("Manifest: %+v\n", manifest)
	fmt.Printf("Image Config: %+v\n", config)

}
