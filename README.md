# Week 6 - Container Image Inspector
A CLI tool written in Go that inspects Docker/OCI container images directly from a registry, without requiring Docker to be installed. It pulls image manifests, lists layers with their sizes, and displays the Dockerfile commands that created each layer.
 
Built as Week 6 of a 12-week project-based learning plan focused on Go and C#.
 
## Features
 
- Authenticates with Docker Hub using the OCI Distribution token flow
- Resolves multi-architecture manifest indexes (defaults to linux/amd64)
- Fetches and parses image manifests and configuration blobs
- Correlates layers with their Dockerfile history entries
- Displays a formatted summary with layer sizes and commands
- Optionally inspects individual layer contents (gzip/tar extraction)
## Usage
 
```
go run . <image> <tag> [--inspect-layer N]
```
 
### Examples
 
```
go run . library/alpine latest
go run . library/nginx latest
go run . library/python 3.12
go run . library/alpine latest --inspect-layer 1
```
 
### Sample Output
 
```
Fetching manifest for library/nginx:latest...
 
================================================================================
Image:         library/nginx:latest
OS/Arch:       linux/amd64
Total Size:    60.03 MB
Layers:        7
================================================================================
 
#  SIZE      COMMAND
-  ----      -------
1  28.40 MB  # debian.sh --arch 'amd64' out/ 'trixie' '@1777939200'
             LABEL maintainer=NGINX Docker Maintainers <docker-maint@nginx.com>
             ENV NGINX_VERSION=1.29.8
             ...
2  31.62 MB  RUN /bin/sh -c set -x     && groupadd --system --gid 101 nginx...
3  628 B     COPY docker-entrypoint.sh / # buildkit
4  956 B     COPY 10-listen-on-ipv6-by-default.sh /docker-entrypoint.d # buildkit
5  404 B     COPY 15-local-resolvers.envsh /docker-entrypoint.d # buildkit
6  1.18 KB   COPY 20-envsubst-on-templates.sh /docker-entrypoint.d # buildkit
7  1.37 KB   COPY 30-tune-worker-processes.sh /docker-entrypoint.d # buildkit
             ENTRYPOINT ["/docker-entrypoint.sh"]
             EXPOSE map[80/tcp:{}]
             STOPSIGNAL SIGQUIT
             CMD ["nginx" "-g" "daemon off;"]
```
 
## Project Structure
 
```
container-inspector/
├── cmd/
│   └── inspector/
│       └── main.go              # CLI entry point, argument parsing
├── internal/
│   ├── registry/
│   │   ├── client.go            # HTTP client, OCI auth token flow
│   │   ├── manifest.go          # Manifest and image config fetching
│   │   └── types.go             # Structs for manifests, configs, descriptors
│   ├── image/
│   │   ├── config.go            # Dockerfile command cleanup/parsing
│   │   ├── layer.go             # Layer content inspection (gzip/tar)
│   │   └── summary.go           # Correlates layers with build history
│   └── output/
│       └── formatter.go         # Tabwriter-based formatted output
├── go.mod
└── README.md
```
 
## How It Works
 
The tool interacts directly with the OCI Distribution API without any Docker dependencies:
 
1. **Authentication** — Makes an unauthenticated request to the registry, parses the `WWW-Authenticate` header from the 401 response, and exchanges it for a bearer token from Docker Hub's auth service.
2. **Manifest Resolution** — Fetches the manifest for the given image and tag. If the registry returns a manifest index (multi-architecture list), it selects the linux/amd64 entry and fetches that platform's specific manifest.
3. **Config Retrieval** — Uses the config digest from the manifest to fetch the image configuration blob, which contains the build history with Dockerfile commands.
4. **Layer Correlation** — Maps the manifest's layer array (which has sizes and digests) to the config's history array (which has commands). History entries marked as `empty_layer` (metadata-only instructions like `CMD`, `ENV`, `EXPOSE`) are displayed without a layer number.
5. **Layer Inspection** (optional) — Downloads a specific layer blob, decompresses it with gzip, and reads the tar archive to list all files and directories within that layer.
## Concepts Practiced
 
- OCI Distribution Spec registry HTTP API
- HTTP client work with `net/http`, including auth token flows and header parsing
- JSON parsing with `encoding/json` into nested Go structs
- Gzip decompression and tar archive reading (`compress/gzip`, `archive/tar`)
- Idiomatic Go project structure with `internal/` packages
- Formatted CLI output with `text/tabwriter`
## Debug Configuration
 
The project includes a VS Code `launch.json` configured with default arguments:
 
```json
"args": ["library/alpine", "latest"]
```
 
Adjust the args to inspect different images during development.
