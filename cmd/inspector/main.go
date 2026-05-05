package main

import (
	"container-inspector/internal/registry"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: inspector <image> <tag>")
		fmt.Println("Example: inspector library/alpine latest")
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

	fmt.Printf("Auth Token: %s\n", token)
}
