# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project that generates ASCII art visualization of the Mandelbrot set. The main program creates various views including a full set view and zoomed regions showing interesting fractal patterns.

## Development Commands

### Running the Program
```bash
go run ./mandelbrot_ascii.go
```

### Building
```bash
go build -o mandelbrot mandelbrot_ascii.go
```

### Project Status
The code compiles and runs successfully. All major features have been implemented:
- ✅ Replaced unused "math" import with "strings" package
- ✅ Fixed string multiplication syntax using `strings.Repeat("=", 50)`
- ✅ Added file saving functionality with random names, descriptions, and timestamps
- ✅ Implemented HTTP service with interactive web interface
- ✅ Docker containerization with multi-stage build and security hardening

### File Saving Features
- Automatically saves all generated ASCII art to `mandelbrot_gallery.txt`
- Each entry includes timestamp, randomly generated name, and poetic description
- Content is appended to preserve previous generations
- Clean formatting with separators between entries

### HTTP Service Features
- ✅ Web-based interface for interactive Mandelbrot generation
- Real-time ASCII art generation with AJAX
- Preset buttons for interesting fractal regions
- Customizable parameters (dimensions, iterations, coordinate ranges)
- Gallery browser for viewing saved ASCII art
- Clean, responsive web interface with matrix-style ASCII display

### HTTP Service Usage
**Start the web server:**
```bash
go run ./mandelbrot_ascii.go -server
```

**Custom port:**
```bash
go run ./mandelbrot_ascii.go -server -port=:3000
```

**Endpoints:**
- `GET /` - Interactive web interface
- `GET /generate` - API endpoint with query parameters (width, height, maxiter, xmin, xmax, ymin, ymax)
- `GET /gallery` - View saved ASCII art gallery

**Example API call:**
```bash
curl "http://localhost:8080/generate?width=40&height=20&maxiter=100&xmin=-2&xmax=1&ymin=-1&ymax=1"
```

### Docker Deployment
Run the Mandelbrot generator in a Docker container for easy deployment and isolation.

**Prerequisites:**
- Docker installed and running on your system
- Docker Compose (optional, for easier management)

**Build and run with Docker:**
```bash
# Build the Docker image
docker build -t mandelbrot-server .

# Run the container on port 3000
docker run -p 3000:3000 mandelbrot-server

# Run with persistent gallery storage
docker run -p 3000:3000 -v mandelbrot_data:/app/data mandelbrot-server
```

**Using Docker Compose (recommended):**
```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

**Docker Features:**
- ✅ Multi-stage build for optimized image size (37.7MB)
- ✅ Non-root user execution for security
- ✅ Health checks for container monitoring
- ✅ Persistent volume mounting for gallery storage
- ✅ Automatic server startup on port 3000
- ✅ Production-ready with Docker Compose orchestration

**Container Access:**
- Web interface: http://localhost:3000
- API endpoint: http://localhost:3000/generate
- Gallery: http://localhost:3000/gallery

## Code Architecture

### Core Components
- **Complex struct**: Custom implementation for complex number operations (Add, Multiply, MagnitudeSquared)
- **mandelbrotIterations()**: Core algorithm that determines escape time for points in the complex plane
- **iterToChar()**: Maps iteration counts to ASCII characters for visualization density
- **generateMandelbrot()**: Main rendering function that maps pixel coordinates to complex plane
- **zoomView()**: Helper for generating zoomed views of specific regions

### Visualization Features
- Configurable dimensions (80x40 default)
- Multiple zoom levels for interesting fractal regions
- ASCII character mapping from sparse (' ') to dense ('#%@') based on escape time
- Predefined interesting regions: Seahorse Valley, Spiral Region, Lightning Region

### Program Flow
1. Generates full Mandelbrot set view (-2.5 to 1.0 real, -1.25 to 1.25 imaginary)
2. Creates zoomed views of mathematically interesting regions
3. Displays high-detail view with increased iterations
4. Shows ASCII legend explaining character meanings

