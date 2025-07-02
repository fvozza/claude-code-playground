package main

import (
	"fmt"
	"math/rand"
	"os"
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

func main() {
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
