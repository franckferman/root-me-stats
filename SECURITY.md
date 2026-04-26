# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | :white_check_mark: |
| 1.x.x   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it to:

- **Email**: [security contact]
- **GitHub**: Create a private security advisory

**Please do not open public issues for security vulnerabilities.**

## Security Measures

### Code Security
- **Zero external dependencies** - Only Go stdlib used
- **No eval/exec** - No dynamic code execution
- **Input validation** - All user inputs validated
- **No SQL injection** - No database, only HTTP requests
- **Path traversal protection** - Only whitelisted file operations

### API Security
- **Rate limiting** - Implemented via caching (24h TTL)
- **CORS configured** - Proper origin handling
- **No sensitive data** - Public Root-me data only
- **No authentication** - Read-only public API

### Deployment Security
- **Single binary** - No complex dependencies
- **Non-root user** - Can run as unprivileged user
- **Minimal attack surface** - HTTP server only
- **No file uploads** - Read-only operations only

### Network Security
- **HTTPS ready** - Works behind reverse proxy
- **No outbound except Root-me** - Only contacts root-me.org
- **Timeout controls** - All HTTP requests have timeouts
- **User-Agent set** - Identifies itself properly

## Best Practices

### Production Deployment
```bash
# Run as non-root user
useradd -r -s /bin/false rootme
su rootme -c './rootme-server'

# Behind reverse proxy (nginx/caddy)
# Enable rate limiting at proxy level
# Use HTTPS termination
```

### Environment Variables
```bash
HOST=127.0.0.1  # Bind to localhost behind proxy
PORT=3000       # Non-privileged port
```

### Monitoring
- Monitor `/health` endpoint
- Watch for excessive requests to root-me.org
- Monitor memory usage (should stay under 50MB)

## Update Policy

Security updates will be released as patch versions and announced via:
- GitHub Releases
- Security advisories (if applicable)

## Acknowledgments

Report security issues responsibly and we'll acknowledge your contribution.