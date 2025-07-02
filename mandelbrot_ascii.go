package main

import (
	"fmt"
	"math"
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

// Generate Mandelbrot ASCII art
func generateMandelbrot(width, height, maxIter int, xMin, xMax, yMin, yMax float64) {
	fmt.Printf("Mandelbrot Set ASCII Art (%dx%d)\n", width, height)
	fmt.Printf("Range: x[%.2f, %.2f], y[%.2f, %.2f]\n", xMin, xMax, yMin, yMax)
	fmt.Printf("Max iterations: %d\n\n", maxIter)
	
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			// Map pixel coordinates to complex plane
			x := xMin + float64(col)*(xMax-xMin)/float64(width-1)
			y := yMin + float64(row)*(yMax-yMin)/float64(height-1)
			
			c := Complex{x, y}
			iter := mandelbrotIterations(c, maxIter)
			char := iterToChar(iter, maxIter)
			
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

// Zoom into an interesting region of the Mandelbrot set
func zoomView(width, height, maxIter int, centerX, centerY, zoom float64) {
	size := 2.0 / zoom
	xMin := centerX - size/2
	xMax := centerX + size/2
	yMin := centerY - size/2
	yMax := centerY + size/2
	
	fmt.Printf("\nZoomed view (zoom: %.1fx)\n", zoom)
	generateMandelbrot(width, height, maxIter, xMin, xMax, yMin, yMax)
}

func main() {
	// Configuration
	width := 80
	height := 40
	maxIter := 100
	
	// Classic full view of Mandelbrot set
	fmt.Println("=== FULL MANDELBROT SET ===")
	generateMandelbrot(width, height, maxIter, -2.5, 1.0, -1.25, 1.25)
	
	// Zoomed views of interesting regions
	fmt.Println("\n" + "="*50)
	fmt.Println("=== INTERESTING REGIONS ===")
	
	// Seahorse valley
	fmt.Println("\n--- Seahorse Valley ---")
	zoomView(width, height, maxIter, -0.75, 0.1, 20)
	
	// Spiral region
	fmt.Println("\n--- Spiral Region ---")
	zoomView(width, height, maxIter, -0.16, 1.04, 50)
	
	// Lightning region
	fmt.Println("\n--- Lightning Region ---")
	zoomView(width, height, maxIter, -1.25, 0.02, 100)
	
	// High detail view with more iterations
	fmt.Println("\n--- High Detail View ---")
	generateMandelbrot(width, height, 200, -0.8, -0.7, 0.05, 0.15)
	
	// ASCII art legend
	fmt.Println("\n" + "="*50)
	fmt.Println("ASCII LEGEND:")
	fmt.Println("' ' (space) - Quick escape (not in set)")
	fmt.Println(".:- - Fast escape")
	fmt.Println("=+* - Medium escape time") 
	fmt.Println("#%@ - Slow escape / In the set")
	fmt.Println("\nThe darker the character, the more iterations")
	fmt.Println("it took to determine if the point escapes.")
}
