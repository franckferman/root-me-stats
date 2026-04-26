# Root-me Stats

Generate Root-me statistics badges for GitHub profiles.

Rewrite of the original Node.js version in Go for better performance and simpler deployment.

## Features

- SVG badge generation with 6 themes
- Profile comparison between users
- CLI tools and HTTP API
- 24h caching to avoid spamming Root-me
- Single binaries with no runtime dependencies
- Cross-platform builds (Linux, Windows, macOS)

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and automated releases:

- **Automated testing** on every push and pull request
- **Cross-platform builds** for Linux, Windows, and macOS (Intel + ARM)
- **Integration tests** verify CLI functionality and server endpoints
- **Automatic releases** when version tags are pushed
- **Pre-compiled binaries** available in GitHub releases

### Build Status

All builds are automatically tested across multiple platforms. Check the Actions tab for current build status.

### Download Pre-built Binaries

Get the latest release from [GitHub Releases](https://github.com/franckferman/root-me-stats/releases) - no compilation needed.

## Usage

Build from source:
```bash
git clone https://github.com/franckferman/root-me-stats.git
cd root-me-stats
make build
```

Generate a badge:
```bash
./bin/rootme-badges --nickname=franckferman --theme=dark --output=badge.svg
```

Start the API server:
```bash
./bin/rootme-server
# Available at http://localhost:3000/rm-gh?nickname=franckferman&style=dark
```

For GitHub Actions (using pre-built binary):
```yaml
- run: |
    curl -L https://github.com/franckferman/root-me-stats/releases/latest/download/rootme-stats-linux-amd64.tar.gz | tar xz
    ./rootme-cli badge --nickname=${{ github.repository_owner }} --output=badge.svg
```

## API Endpoints

### Badge Generation

**Endpoint:** `GET /rm-gh` (compatible) or `GET /badge`

**Parameters:**
- `nickname` (required) - Root-me username
- `style` (optional) - Theme: `dark`, `light`, `midnight`, `punk`, `weedy`, `astral`  
- `gstats` (optional) - Show global stats: `show` or `hidden`

**Example:**
```
https://your-server.com/rm-gh?nickname=franckferman&style=midnight&gstats=show
```

### Profile Comparison

**Endpoint:** `GET /compare`

**Parameters:**
- `user1` (required) - First username
- `user2` (required) - Second username
- `style` (optional) - Theme for comparison
- `width` (optional) - Badge width

### JSON APIs

**Profile data:** `GET /api/profile?nickname=USERNAME`  
**Comparison data:** `GET /api/compare?user1=USER1&user2=USER2`

## CLI Usage

### Generate Badge

```bash
# Basic usage
./rootme-cli badge --nickname=franckferman --output=./badge.svg

# With options
./rootme-cli badge \
  --nickname=franckferman \
  --output=./assets/rootme-badge.svg \
  --theme=dark \
  --stats
```

### Generate Comparison

```bash
./rootme-cli compare \
  --user1=user1 \
  --user2=user2 \
  --output=./compare.svg \
  --theme=midnight
```

### Fetch Profile Data

```bash
# Save to file
./rootme-cli profile --nickname=franckferman --output=./profile.json

# Output to stdout (for GitHub Actions)
./rootme-cli profile --nickname=franckferman
```

### Simple Badge Tool

```bash
# Lightweight tool for basic badges
./rootme-badges --nickname=franckferman --theme=dark --stats
./rootme-badges --nickname=franckferman --output=badge.svg --theme=midnight
```

## GitHub Actions Integration

### Automatic Badge Updates

Create `.github/workflows/update-rootme-badge.yml`:

```yaml
name: Update Root-me Badge

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  workflow_dispatch:     # Manual trigger

jobs:
  update-badge:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Download Root-me Stats
        run: |
          curl -L https://github.com/franckferman/root-me-stats/releases/latest/download/rootme-stats-linux-amd64.tar.gz | tar xz

      - name: Generate Root-me Badge  
        run: |
          ./rootme-cli badge \
            --nickname=${{ github.repository_owner }} \
            --output=./assets/rootme-badge.svg \
            --theme=dark \
            --stats

      - name: Commit and push badge
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add assets/rootme-badge.svg
          git diff --staged --quiet || git commit -m "🔄 Update Root-me badge"
          git push
```

### Use in README

```markdown
<!-- In your GitHub profile README.md -->
<p align="center">
  <img src="https://raw.githubusercontent.com/USERNAME/USERNAME/main/assets/rootme-badge.svg" alt="Root-me Stats" />
</p>

<!-- Or reference your deployed service -->
<p align="center">
  <img src="https://your-server.com/rm-gh?nickname=USERNAME&style=dark&gstats=show" alt="Root-me Stats" />
</p>
```

## Deployment

**Local/VPS (recommended):**
```bash
./rootme-server  # Just run the binary
# Or: systemctl service, Docker, whatever you prefer
```

**Cloud platforms:**
- **Vercel/Netlify:** Deploy binary as function
- **Heroku/Railway:** Direct binary deployment  
- **VPS/Dedicated:** Copy file and run

**Cross-platform builds:**
```bash
make build-all  # Linux, Windows, macOS (Intel + ARM)
```

## Examples

### Badge Themes

Available themes: `dark`, `light`, `midnight`, `punk`, `weedy`, `astral`

### Profile Comparison

Compare two Root-me profiles side by side using the `/compare` endpoint.

## Go Library Usage

```go
import "github.com/franckferman/root-me-stats/pkg/rootme"

// Quick badge generation
svg, err := rootme.QuickBadge("franckferman", rootme.BadgeOptions{
    Theme: "midnight", 
    ShowGlobalStats: true,
})

// Profile data
profile, err := rootme.GetProfile("franckferman")

// Comparison
comparison, err := rootme.CompareProfiles("user1", "user2")
```

## Architecture

```
cmd/
├── server/     # HTTP API server (single binary)
├── cli/        # Full CLI tool with all features  
└── badges/     # Lightweight badge-only tool

internal/
├── fetcher/    # Root-me data extraction (stdlib only)
├── generator/  # SVG badge generation
├── themes/     # Theme definitions
└── cache/      # Simple file-based cache

pkg/
└── rootme/     # Public Go API
```

**Design principles:**
- **Zero external dependencies** - Go stdlib only
- **Single binaries** - No runtime dependencies
- **Fast and secure** - Minimal attack surface
- **GitHub Actions native** - CLI-first design  
- **Cross-platform** - Builds everywhere Go runs

## Why Go?

The original Node.js version had some issues:
- Large runtime dependencies (node_modules)
- Slow cold start times
- Complex deployment

This Go version fixes those problems:
- Single binary with no dependencies
- Fast startup and low memory usage
- Cross-compilation for different platforms

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`  
3. Build and test: `make build && make test-real`
4. Commit changes: `git commit -m 'Add amazing feature'`
5. Push and open PR: `git push origin feature/amazing-feature`

## License

MIT License - feel free to use in your projects!

## Credits

Inspired by and compatible with:
- [Rootme-readme-stats](https://github.com/dz-root/Rootme-readme-stats) by dz-root
- [Root-me-diff](https://github.com/dz-root/Root-me-diff) by dz-root

**Now in Go for maximum performance and zero dependencies!** 🚀