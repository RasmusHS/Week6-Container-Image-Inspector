package main

import (
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
}
