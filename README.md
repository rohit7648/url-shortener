# URL Shortener Service

A high-performance URL shortener service built with Go and Kratos framework.

## Features

- Shorten long URLs to manageable short URLs
- Redirect short URLs to original URLs
- RESTful API endpoints
- gRPC support
- Configuration management
- Docker support

## Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose (optional)
- Make
- Protocol Buffers compiler (protoc)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/rohit7648/url-shortener.git
cd url-shortener
```

2. Install development dependencies:
```bash
make init
```

3. Generate protobuf files:
```bash
make all
```

## Configuration

The service uses a YAML configuration file located at `configs/config.yaml`. You can modify this file to change the service settings.

## Running the Service

### Development Mode

To run the service in development mode with hot reload:
```bash
make dev
```

### Production Mode

To build and run the service:
```bash
make run
```

### Using Docker

1. Build the Docker image:
```bash
docker build -t url-shortener .
```

2. Run the container:
```bash
docker run -p 8000:8000 -p 9000:9000 url-shortener
```

## API Documentation

### REST API

The service exposes the following REST endpoints:

- `POST /v1/url/shorten` - Create a short URL
- `GET /v1/url/{short_url}` - Get original URL
- `GET /v1/url/{short_url}/redirect` - Redirect to original URL

### gRPC API

The service also provides gRPC endpoints with the same functionality.

## Available Make Commands

- `make init` - Install development dependencies
- `make config` - Generate internal protobuf files
- `make api` - Generate API protobuf files
- `make build` - Build the project
- `make generate` - Generate wire and other files
- `make all` - Generate all protobuf files and wire files
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make dev` - Run in development mode with hot reload
- `make run` - Build and run the application

## Project Structure

```
.
├── api/                    # API definitions
│   └── urlshortener/      # URL shortener API
├── cmd/                    # Main applications
├── configs/               # Configuration files
├── internal/              # Private application code
│   ├── biz/              # Business logic
│   ├── conf/             # Configuration
│   ├── data/             # Data layer
│   ├── server/           # Server implementations
│   └── service/          # Service implementations
├── third_party/          # Third-party dependencies
└── Makefile              # Build and development commands
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


