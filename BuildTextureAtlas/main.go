package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load individual images and pack into a square (power of 2) larger image.

// https://codeincomplete.com/articles/bin-packing/
// http://code.activestate.com/recipes/442299/
// https://gamedev.stackexchange.com/questions/2829/texture-packing-algorithm
// https://www.cs.upc.edu/~jmartinez/slides/masterThesisSlides.pdf

// Collect all the input items.
// Sort them by total pixels consumed, large-to-small. Prioritize on Squares
// over Rectangles. Lay them out in your
// texture in scanline order, just testing stuff from the top-left pixel to
// the top-right pixel, moving down a line, and repeating, resetting to the
// top-left pixel after every successful placement.

// You either need to hardcode a width or come up with another heuristic
// for this. In an attempt to preserve squareness, our algorithm would
// start at 128, then increase by 128s until it came up with a result that
// wasn't any deeper than it was wide.

// Start with 128x128 and attempt to fit all images. If not then bump
// up to 256x256 etc.

type cell struct {
	image  image.Image
	w, h   int
	area   int
	placed bool
	square bool
	// coords in atlas
	px, py int
}

func main() {
	testPacker()

	// Open manifest file
	maniLines := loadManifest("manifest.txt")
	cells := []*cell{}

	// ----------------------------------------------------------
	// Load all images
	// ----------------------------------------------------------
	for _, source := range maniLines {
		imageFiles := collectFiles(source)
		for _, imagePath := range imageFiles {
			img, err := loadImage(imagePath)
			if err != nil {
				panic(err)
			}

			cell := &cell{}
			cell.image = img
			cell.w = img.Bounds().Dx()
			cell.h = img.Bounds().Dy()
			cell.area = cell.w * cell.h
			cell.square = cell.w == cell.h
			cells = append(cells, cell)
		}
	}

	fmt.Println("Loaded images")

	// ----------------------------------------------------------
	// Gather stats from images, for example, max/min width and height
	// We need to sort the images by Area, prioritizing on Square over
	// Rectangle.
	// ----------------------------------------------------------

	for _, cell := range cells {
		if cell.square {
		} else {

		}
	}

}

func collectFiles(path string) []string {
	if strings.Contains(path, "*") {
		// Get file listing of directory
		var files []string
		root := strings.TrimRight(path, "*")

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsRegular() {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

		return files
	}

	return []string{path}
}

func testPacker() {
	// 100x100x2
	// 80x80x2
	// 50x50x2
	// 25x25x3
	// 100x25x4

	packSize := int(math.Pow(2, 9)) // 8 = 256, 9 = 512 etc...

	blocks := []*node{
		{w: 128, h: 128},
		{w: 128, h: 128},
		{w: 128, h: 32},
		{w: 128, h: 32},
		{w: 128, h: 32},
		{w: 128, h: 32},
		{w: 128, h: 128},
		{w: 128, h: 128},
		{w: 64, h: 64},
		{w: 64, h: 64},
		{w: 32, h: 32},
		{w: 32, h: 32},
		{w: 32, h: 32},
	}

	// This sorts by height. Even better would be to then sort by width
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[j].h < blocks[i].h
	})

	// for _, block := range blocks {
	// 	fmt.Println(block)
	// }

	pac := packer{root: &node{w: packSize, h: packSize}}

	pac.fit(blocks)

	var maxX, maxY int

	packedBlocks := 0
	fitTotal := 0

	for _, block := range blocks {
		if block.fit != nil {
			fit := block.fit
			fmt.Printf("x:%03d, y: %03d, w: %03d, h: %03d, area: %d\n", fit.x, fit.y, fit.w, fit.h, fit.area)
			maxX = int(math.Max(float64(maxX), float64(fit.x)))
			maxY = int(math.Max(float64(maxY), float64(fit.y)))
			packedBlocks++
			fitTotal += fit.area
		}
	}

	fmt.Printf("maxX: %d, maxY: %d\n", maxX, maxY)

	totalBlocks := len(blocks)

	fmt.Println("Container size: ", packSize)
	fmt.Println("Total blocks: ", totalBlocks)
	fmt.Println("Total packed blocks: ", packedBlocks)
	if totalBlocks > packedBlocks {
		fmt.Println("##################################################")
		fmt.Println("*** Warning!! Not all blocks could be packed!! ***")
		fmt.Println("##################################################")
	} else {
		fmt.Printf("Filled: %d%%\n", int(math.Round(100.0*float64(fitTotal)/float64(packSize*packSize))))
	}

	fmt.Println("Done")
}
