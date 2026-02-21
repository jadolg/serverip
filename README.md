# Server IP Details (serverip)

A simple, lightweight Go-based utility to identify the public exit IP address of a server. Useful for debugging network configurations such as VPNs or tunnels.

Data is provided by [wtfismyip.com](https://wtfismyip.com).

## Features

-   **IP Discovery:** Shows IP address, location, ISP, and Tor exit status.
-   **Web Interface:** Clean, responsive UI with dark mode support.
-   **JSON API:** Detects the `Accept: application/json` header for scriptable access.
-   **Health Checks:** Includes a `/health` endpoint for container orchestration.
-   **Configurable:** Port can be set via the `PORT` environment variable (defaults to 8080).
-   **Robust:** Uses structured logging (`slog`) and timeout-controlled HTTP client.

## Running

### Local

```bash
go run .
```

Visit `http://localhost:8080` in your browser.

### Docker

```bash
docker build -t serverip .
docker run -p 8080:8080 serverip
```

## API Access

To get raw JSON data (e.g., from a CLI):

```bash
curl -H "Accept: application/json" http://localhost:8080
```

## Development

Run tests:

```bash
go test -v ./...
```
