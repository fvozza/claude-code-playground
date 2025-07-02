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
The code compiles and runs successfully. Previous compilation errors have been resolved:
- ✅ Replaced unused "math" import with "strings" package
- ✅ Fixed string multiplication syntax using `strings.Repeat("=", 50)`

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

