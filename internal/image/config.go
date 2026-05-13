package image

import "strings"

// CleanCommand takes a raw created_by string from image history and returns
// a human-readable Dockerfile command
//
// Raw history entries look like:
//   "/bin/sh -c #(nop)  CMD [\"/bin/sh\"]"   → "CMD [\"/bin/sh\"]"
//   "/bin/sh -c #(nop) ADD file:abc in / "    → "ADD file:abc in /"
//   "/bin/sh -c apk add --no-cache curl"      → "RUN apk add --no-cache curl"
func CleanCommand(raw string) string {
	// Handle #(nop) entries — these are metadata instructions (CMD, ENV, EXPOSE, etc.)
	if strings.Contains(raw, "#(nop)") {
		parts := strings.SplitN(raw, "#(nop)", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1])
		}
	}

	// Handle regular RUN commands — these are "/bin/sh -c <actual command>"
	prefix := "/bin/sh -c "
	if strings.HasPrefix(raw, prefix) {
		return "RUN " + strings.TrimPrefix(raw, prefix)
	}

	// Fallback: return as-is
	return raw
}
