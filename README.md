# Customize

A simple service that lets you create customizable API endpoints with configurable behaviors like delays, status codes, and response bodies.

## Running locally

```bash
docker run -it --rm -p 8080:8080 $(docker build -q .)
```

## CORS Configuration

You can configure CORS (Cross-Origin Resource Sharing) headers in the response by setting the following environment variables:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| `CORS_ENABLED` | Enable CORS headers | `false` |
| `CORS_ALLOW_ORIGIN` | Allowed origins | `*` (all origins) |
| `CORS_ALLOW_METHODS` | Allowed HTTP methods | `GET, POST, PUT, DELETE, OPTIONS` |
| `CORS_ALLOW_HEADERS` | Allowed headers | `Content-Type, Authorization` |
| `CORS_MAX_AGE` | How long the preflight response can be cached (in seconds) | `86400` (24 hours) |
| `CORS_ALLOW_CREDENTIALS` | Allow credentials | Not set by default |
| `CORS_EXPOSE_HEADERS` | Headers that can be exposed to the client | Not set by default |

### Example Usage

To enable CORS with default settings:

```bash
docker run -it --rm -p 8080:8080 -e CORS_ENABLED=true $(docker build -q .)
```

To enable CORS with custom settings:

```bash
docker run -it --rm -p 8080:8080 \
  -e CORS_ENABLED=true \
  -e CORS_ALLOW_ORIGIN="https://example.com" \
  -e CORS_ALLOW_METHODS="GET, POST" \
  -e CORS_ALLOW_HEADERS="Content-Type, X-Custom-Header" \
  -e CORS_MAX_AGE="3600" \
  -e CORS_ALLOW_CREDENTIALS="true" \
  -e CORS_EXPOSE_HEADERS="X-Custom-Header" \
  $(docker build -q .)
```
