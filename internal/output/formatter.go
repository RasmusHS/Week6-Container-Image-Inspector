package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"container-inspector/internal/image"
)

// FormatSize converts bytes to a human-readable string (B, KB, MB, GB)
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	// Use a switch statement to determine the appropriate unit for the given byte size and format it accordingly
	// If the byte size is greater than or equal to 1 GB, format it as GB; if it's greater than or equal to 1 MB, format it as MB;
	// if it's greater than or equal to 1 KB, format it as KB;
	// otherwise, format it as bytes (B)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// PrintSummary outputs the image analysis in a formatted table
func PrintSummary(summary *image.Summary) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Image:         %s:%s\n", summary.Image, summary.Tag)
	fmt.Printf("OS/Arch:       %s/%s\n", summary.OS, summary.Architecture)
	fmt.Printf("Total Size:    %s\n", FormatSize(summary.TotalSize))
	fmt.Printf("Layers:        %d\n", countRealLayers(summary))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Create a new tab writer to format the output in a table-like structure with columns for layer index, size, and command
	// The header row includes the column titles: "#", "SIZE", and "COMMAND"
	// We also print a separator row with dashes to visually separate the header from the layer entries
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tSIZE\tCOMMAND")
	fmt.Fprintln(w, "-\t----\t-------")

	// Loop through the layers in the summary and print their information in the table
	// For each layer, we clean the created_by command for better readability and truncate it if it's too long
	// If the layer is an empty layer, we print a placeholder in the index and size columns and only show the command; otherwise, we print the actual index, formatted size, and command
	for _, layer := range summary.Layers {
		cmd := image.CleanCommand(layer.CreatedBy)

		// Truncate long commands for display
		if len(cmd) > 80 {
			cmd = cmd[:77] + "..."
		}

		if layer.EmptyLayer {
			fmt.Fprintf(w, " \t \t%s\n", cmd)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n", layer.Index, FormatSize(layer.Size), cmd)
		}
	}

	// Flush the tab writer to ensure all output is written to the console, and print a final newline for spacing
	w.Flush()
	fmt.Println()
}

// PrintLayerContents outputs the file listing for an inspected layer
func PrintLayerContents(layerIndex int, digest string, entries []image.LayerEntry) {
	fmt.Printf("\nLayer %d contents (%s):\n", layerIndex, digest[:19])
	fmt.Println(strings.Repeat("-", 60))

	// Create a new tab writer to format the output in a table-like structure with columns for type, size, and name
	// The header row includes the column titles: "TYPE", "SIZE", and "NAME"
	// We also print a separator row with dashes to visually separate the header from the file entries
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tSIZE\tNAME")

	// Loop through the file entries in the layer and print their information in the table
	// For each entry, we print the type (file, directory, symlink), formatted size, and name of the file or directory
	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", entry.Type, FormatSize(entry.Size), entry.Name)
	}

	// Flush the tab writer to ensure all output is written to the console, and print a final newline for spacing
	w.Flush()
	fmt.Printf("\nTotal entries: %d\n", len(entries))
}

// countRealLayers counts the number of non-empty layers in the summary by iterating through the layers and incrementing a counter for each layer that is not marked as an empty layer
func countRealLayers(summary *image.Summary) int {
	count := 0
	for _, l := range summary.Layers {
		if !l.EmptyLayer {
			count++
		}
	}
	return count
}
