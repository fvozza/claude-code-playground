package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Complex number structure
type Complex struct {
	real, imag float64
}

// Add two complex numbers
func (c Complex) Add(other Complex) Complex {
	return Complex{c.real + other.real, c.imag + other.imag}
}

// Multiply two complex numbers
func (c Complex) Multiply(other Complex) Complex {
	return Complex{
		c.real*other.real - c.imag*other.imag,
		c.real*other.imag + c.imag*other.real,
	}
}

// Calculate the magnitude squared of a complex number
func (c Complex) MagnitudeSquared() float64 {
	return c.real*c.real + c.imag*c.imag
}

// Mandelbrot iteration function
func mandelbrotIterations(c Complex, maxIter int) int {
	z := Complex{0, 0}

	for i := 0; i < maxIter; i++ {
		if z.MagnitudeSquared() > 4.0 {
			return i
		}
		z = z.Multiply(z).Add(c)
	}
	return maxIter
}

// Convert iteration count to ASCII character
func iterToChar(iter, maxIter int) rune {
	chars := []rune(" .:-=+*#%@")
	if iter == maxIter {
		return chars[len(chars)-1] // Inside the set
	}

	// Map iteration to character index
	index := int(float64(iter) / float64(maxIter) * float64(len(chars)-1))
	if index >= len(chars) {
		index = len(chars) - 1
	}
	return chars[index]
}

// Generate a random name and description for the ASCII art
func generateArtMetadata() (string, string) {
	adjectives := []string{"Mystical", "Ethereal", "Cosmic", "Infinite", "Swirling", "Fractal", "Chaotic", "Beautiful", "Complex", "Mathematical"}
	nouns := []string{"Spiral", "Vortex", "Pattern", "Dream", "Universe", "Landscape", "Vision", "Gateway", "Portal", "Dimension"}
	descriptors := []string{"dancing through infinite complexity", "revealing hidden mathematical beauty", "emerging from chaos", "spiraling into eternity", "whispering secrets of infinity", "mapping the edge of existence", "painting mathematics with ASCII", "bridging reality and abstraction", "showing the art within algorithms", "exploring the fractal frontier"}

	rand.Seed(time.Now().UnixNano())

	name := adjectives[rand.Intn(len(adjectives))] + " " + nouns[rand.Intn(len(nouns))]
	description := descriptors[rand.Intn(len(descriptors))]

	return name, description
}

// Save ASCII art to file with metadata
func saveAsciiArt(content, title string) error {
	filename := "mandelbrot_gallery.txt"

	// Check if file exists to determine if we need header
	var file *os.File
	var err error

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, create it with header
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		_, err = file.WriteString("MANDELBROT ASCII ART GALLERY\n")
		if err != nil {
			file.Close()
			return err
		}
		_, err = file.WriteString(strings.Repeat("=", 80) + "\n\n")
		if err != nil {
			file.Close()
			return err
		}
	} else {
		// File exists, open for appending
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}

	defer file.Close()

	// Generate metadata
	name, description := generateArtMetadata()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Write separator and metadata
	_, err = file.WriteString("\n" + strings.Repeat("-", 80) + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("Generated: %s\n", timestamp))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("Name: %s\n", name))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("Description: %s\n", description))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("Type: %s\n\n", title))
	if err != nil {
		return err
	}

	// Write the ASCII art content
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n\n")
	if err != nil {
		return err
	}

	return nil
}

// MandelbrotParams holds parameters for generating Mandelbrot art
type MandelbrotParams struct {
	Width   int
	Height  int
	MaxIter int
	XMin    float64
	XMax    float64
	YMin    float64
	YMax    float64
	Title   string
}

// Preset regions for interesting Mandelbrot views
var presets = map[string]MandelbrotParams{
	"full": {
		Width: 80, Height: 40, MaxIter: 100,
		XMin: -2.5, XMax: 1.0, YMin: -1.25, YMax: 1.25,
		Title: "Full Mandelbrot Set",
	},
	"seahorse": {
		Width: 80, Height: 40, MaxIter: 100,
		XMin: -0.8, XMax: -0.7, YMin: 0.05, YMax: 0.15,
		Title: "Seahorse Valley",
	},
	"spiral": {
		Width: 80, Height: 40, MaxIter: 100,
		XMin: -0.18, XMax: -0.14, YMin: 1.02, YMax: 1.06,
		Title: "Spiral Region",
	},
	"lightning": {
		Width: 80, Height: 40, MaxIter: 100,
		XMin: -1.26, XMax: -1.24, YMin: 0.01, YMax: 0.03,
		Title: "Lightning Region",
	},
	"highdetail": {
		Width: 80, Height: 40, MaxIter: 200,
		XMin: -0.8, XMax: -0.7, YMin: 0.05, YMax: 0.15,
		Title: "High Detail View",
	},
}

// Generate Mandelbrot ASCII art using parameters
func generateMandelbrotWithParams(params MandelbrotParams) string {
	return generateMandelbrot(params.Width, params.Height, params.MaxIter, params.XMin, params.XMax, params.YMin, params.YMax)
}

// Generate Mandelbrot ASCII art to an io.Writer (for HTTP responses)
func generateMandelbrotToWriter(w io.Writer, params MandelbrotParams, includeHeader bool) {
	if includeHeader {
		header := fmt.Sprintf("Mandelbrot Set ASCII Art (%dx%d)\n", params.Width, params.Height)
		header += fmt.Sprintf("Range: x[%.2f, %.2f], y[%.2f, %.2f]\n", params.XMin, params.XMax, params.YMin, params.YMax)
		header += fmt.Sprintf("Max iterations: %d\n\n", params.MaxIter)
		fmt.Fprint(w, header)
	}

	for row := 0; row < params.Height; row++ {
		for col := 0; col < params.Width; col++ {
			// Map pixel coordinates to complex plane
			x := params.XMin + float64(col)*(params.XMax-params.XMin)/float64(params.Width-1)
			y := params.YMin + float64(row)*(params.YMax-params.YMin)/float64(params.Height-1)

			c := Complex{x, y}
			iter := mandelbrotIterations(c, params.MaxIter)
			char := iterToChar(iter, params.MaxIter)

			fmt.Fprint(w, string(char))
		}
		fmt.Fprint(w, "\n")
	}
}

// Generate Mandelbrot ASCII art
func generateMandelbrot(width, height, maxIter int, xMin, xMax, yMin, yMax float64) string {
	var result strings.Builder

	header := fmt.Sprintf("Mandelbrot Set ASCII Art (%dx%d)\n", width, height)
	header += fmt.Sprintf("Range: x[%.2f, %.2f], y[%.2f, %.2f]\n", xMin, xMax, yMin, yMax)
	header += fmt.Sprintf("Max iterations: %d\n\n", maxIter)

	fmt.Print(header)
	result.WriteString(header)

	for row := 0; row < height; row++ {
		var line strings.Builder
		for col := 0; col < width; col++ {
			// Map pixel coordinates to complex plane
			x := xMin + float64(col)*(xMax-xMin)/float64(width-1)
			y := yMin + float64(row)*(yMax-yMin)/float64(height-1)

			c := Complex{x, y}
			iter := mandelbrotIterations(c, maxIter)
			char := iterToChar(iter, maxIter)

			line.WriteRune(char)
		}
		lineStr := line.String()
		fmt.Println(lineStr)
		result.WriteString(lineStr + "\n")
	}

	return result.String()
}

// Zoom into an interesting region of the Mandelbrot set
func zoomView(width, height, maxIter int, centerX, centerY, zoom float64) string {
	size := 2.0 / zoom
	xMin := centerX - size/2
	xMax := centerX + size/2
	yMin := centerY - size/2
	yMax := centerY + size/2

	fmt.Printf("\nZoomed view (zoom: %.1fx)\n", zoom)
	art := generateMandelbrot(width, height, maxIter, xMin, xMax, yMin, yMax)
	return fmt.Sprintf("Zoomed view (zoom: %.1fx)\n%s", zoom, art)
}

// HTTP handlers

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Mandelbrot ASCII Art Generator</title>
    <style>
        body { font-family: monospace; margin: 40px; background: #f0f0f0; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; }
        .controls { background: #f8f8f8; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .preset-buttons { margin: 10px 0; }
        .preset-buttons button { margin: 5px; padding: 10px 15px; cursor: pointer; }
        .ascii-output { background: #000; color: #0f0; padding: 20px; border-radius: 5px; overflow-x: auto; }
        pre { margin: 0; font-size: 12px; line-height: 1.2; }
        input, select { margin: 5px; }
        .form-group { margin: 10px 0; }
        .form-group label { display: inline-block; width: 120px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Mandelbrot ASCII Art Generator</h1>
        
        <div class="controls">
            <h3>Quick Presets</h3>
            <div class="preset-buttons">
                <button onclick="loadPreset('full')">Full Set</button>
                <button onclick="loadPreset('seahorse')">Seahorse Valley</button>
                <button onclick="loadPreset('spiral')">Spiral Region</button>
                <button onclick="loadPreset('lightning')">Lightning Region</button>
                <button onclick="loadPreset('highdetail')">High Detail</button>
            </div>
            
            <h3>Custom Parameters</h3>
            <form id="mandelbrotForm">
                <div class="form-group">
                    <label>Width:</label>
                    <input type="number" name="width" value="80" min="10" max="200">
                </div>
                <div class="form-group">
                    <label>Height:</label>
                    <input type="number" name="height" value="40" min="10" max="100">
                </div>
                <div class="form-group">
                    <label>Max Iterations:</label>
                    <input type="number" name="maxiter" value="100" min="10" max="1000">
                </div>
                <div class="form-group">
                    <label>X Min:</label>
                    <input type="number" name="xmin" value="-2.5" step="0.01">
                </div>
                <div class="form-group">
                    <label>X Max:</label>
                    <input type="number" name="xmax" value="1.0" step="0.01">
                </div>
                <div class="form-group">
                    <label>Y Min:</label>
                    <input type="number" name="ymin" value="-1.25" step="0.01">
                </div>
                <div class="form-group">
                    <label>Y Max:</label>
                    <input type="number" name="ymax" value="1.25" step="0.01">
                </div>
                <button type="submit">Generate ASCII Art</button>
            </form>
        </div>
        
        <div class="ascii-output">
            <pre id="output">Click "Generate ASCII Art" or choose a preset to see the fractal!</pre>
        </div>
        
        <p><a href="/gallery">View Gallery</a></p>
    </div>
    
    <script>
        const presets = {
            'full': {width: 80, height: 40, maxiter: 100, xmin: -2.5, xmax: 1.0, ymin: -1.25, ymax: 1.25},
            'seahorse': {width: 80, height: 40, maxiter: 100, xmin: -0.8, xmax: -0.7, ymin: 0.05, ymax: 0.15},
            'spiral': {width: 80, height: 40, maxiter: 100, xmin: -0.18, xmax: -0.14, ymin: 1.02, ymax: 1.06},
            'lightning': {width: 80, height: 40, maxiter: 100, xmin: -1.26, xmax: -1.24, ymin: 0.01, ymax: 0.03},
            'highdetail': {width: 80, height: 40, maxiter: 200, xmin: -0.8, xmax: -0.7, ymin: 0.05, ymax: 0.15}
        };
        
        function loadPreset(name) {
            const preset = presets[name];
            const form = document.getElementById('mandelbrotForm');
            form.width.value = preset.width;
            form.height.value = preset.height;
            form.maxiter.value = preset.maxiter;
            form.xmin.value = preset.xmin;
            form.xmax.value = preset.xmax;
            form.ymin.value = preset.ymin;
            form.ymax.value = preset.ymax;
            generateArt();
        }
        
        function generateArt() {
            const form = document.getElementById('mandelbrotForm');
            const formData = new FormData(form);
            const params = new URLSearchParams(formData);
            
            document.getElementById('output').textContent = 'Generating...';
            
            fetch('/generate?' + params.toString())
                .then(response => response.text())
                .then(data => {
                    document.getElementById('output').textContent = data;
                })
                .catch(error => {
                    document.getElementById('output').textContent = 'Error: ' + error;
                });
        }
        
        document.getElementById('mandelbrotForm').addEventListener('submit', function(e) {
            e.preventDefault();
            generateArt();
        });
        
        // Load full set by default
        loadPreset('full');
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse parameters
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	if width == 0 {
		width = 80
	}

	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	if height == 0 {
		height = 40
	}

	maxIter, _ := strconv.Atoi(r.URL.Query().Get("maxiter"))
	if maxIter == 0 {
		maxIter = 100
	}

	xMin, _ := strconv.ParseFloat(r.URL.Query().Get("xmin"), 64)
	if xMin == 0 {
		xMin = -2.5
	}

	xMax, _ := strconv.ParseFloat(r.URL.Query().Get("xmax"), 64)
	if xMax == 0 {
		xMax = 1.0
	}

	yMin, _ := strconv.ParseFloat(r.URL.Query().Get("ymin"), 64)
	if yMin == 0 {
		yMin = -1.25
	}

	yMax, _ := strconv.ParseFloat(r.URL.Query().Get("ymax"), 64)
	if yMax == 0 {
		yMax = 1.25
	}

	params := MandelbrotParams{
		Width: width, Height: height, MaxIter: maxIter,
		XMin: xMin, XMax: xMax, YMin: yMin, YMax: yMax,
		Title: "Custom Generation",
	}

	w.Header().Set("Content-Type", "text/plain")
	generateMandelbrotToWriter(w, params, true)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("mandelbrot_gallery.txt")
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Mandelbrot Gallery</title>
    <style>
        body { font-family: monospace; margin: 40px; background: #f0f0f0; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; }
        .gallery-content { background: #000; color: #0f0; padding: 20px; border-radius: 5px; overflow-x: auto; }
        pre { margin: 0; font-size: 12px; line-height: 1.2; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Mandelbrot ASCII Art Gallery</h1>
        <p><a href="/">&lt; Back to Generator</a></p>
        
        <div class="gallery-content">
            <pre>` + string(content) + `</pre>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func runCLI() {
	// Configuration
	width := 80
	height := 40
	maxIter := 100

	// Classic full view of Mandelbrot set
	fmt.Println("=== FULL MANDELBROT SET ===")
	fullSetArt := generateMandelbrot(width, height, maxIter, -2.5, 1.0, -1.25, 1.25)
	saveAsciiArt(fullSetArt, "Full Mandelbrot Set")

	// Zoomed views of interesting regions
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("=== INTERESTING REGIONS ===")

	// Seahorse valley
	fmt.Println("\n--- Seahorse Valley ---")
	seahorseArt := zoomView(width, height, maxIter, -0.75, 0.1, 20)
	saveAsciiArt(seahorseArt, "Seahorse Valley")

	// Spiral region
	fmt.Println("\n--- Spiral Region ---")
	spiralArt := zoomView(width, height, maxIter, -0.16, 1.04, 50)
	saveAsciiArt(spiralArt, "Spiral Region")

	// Lightning region
	fmt.Println("\n--- Lightning Region ---")
	lightningArt := zoomView(width, height, maxIter, -1.25, 0.02, 100)
	saveAsciiArt(lightningArt, "Lightning Region")

	// High detail view with more iterations
	fmt.Println("\n--- High Detail View ---")
	highDetailArt := generateMandelbrot(width, height, 200, -0.8, -0.7, 0.05, 0.15)
	saveAsciiArt(highDetailArt, "High Detail View")

	// ASCII art legend
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("ASCII LEGEND:")
	fmt.Println("' ' (space) - Quick escape (not in set)")
	fmt.Println(".:- - Fast escape")
	fmt.Println("=+* - Medium escape time")
	fmt.Println("#%@ - Slow escape / In the set")
	fmt.Println("\nThe darker the character, the more iterations")
	fmt.Println("it took to determine if the point escapes.")

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("ASCII art saved to mandelbrot_gallery.txt")
}

func runServer(port string) {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/gallery", galleryHandler)

	fmt.Printf("Starting Mandelbrot HTTP server on http://localhost%s\n", port)
	fmt.Printf("Open your browser and visit: http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	var serverMode = flag.Bool("server", false, "Run as HTTP server instead of CLI")
	var port = flag.String("port", ":8080", "Port for HTTP server (e.g., :8080)")
	flag.Parse()

	if *serverMode {
		runServer(*port)
	} else {
		runCLI()
	}
}
