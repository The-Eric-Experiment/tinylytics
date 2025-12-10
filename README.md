# Tinylytics

A lightweight analytics platform built with Go and HTMX, featuring a retro Windows 95-inspired UI.

## Features

- Real-time analytics tracking
- Session and page view tracking
- Browser, OS, and country detection
- Referrer and page tracking
- SQLite database storage
- Windows 98-style UI using [98.css](https://github.com/jdan/98.css)
- HTMX for dynamic updates without JavaScript frameworks

## Setup

### Prerequisites

- Go 1.24 or higher
- SQLite

### Installation

1. Clone the repository
2. Copy the example config and edit it:

   ```bash
   cp config.yaml.example config.yaml
   # Edit config.yaml with your settings
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8099`

### Docker

Build and run with Docker:

```bash
docker build -t tinylytics .
docker run -p 8099:8099 -v $(pwd)/data:/app/data tinylytics
```

## Configuration

Edit `config.yaml` to configure your websites and data folder:

```yaml
user:
  username: admin
  password: your-password

websites:
  - domain: example.com
    title: Example Website
  - domain: another.com
    title: Another Site

data-folder: ./data
```

## Tracking

Add the tracking script to your website:

```html
<script>
  fetch("https://your-tinylytics-domain.com/api/event", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      domain: "example.com",
      path: window.location.pathname,
      referrer: document.referrer,
    }),
  });
</script>
```

## Development

The application uses:

- **Backend**: Go with Gin framework
- **Database**: SQLite with GORM
- **Frontend**: HTML templates with HTMX
- **Styling**: 98.css for Windows 98 aesthetic + minimal custom CSS

## License

MIT
