package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load individual images and pack into a square (power of 2) larger image.
// Then create a 'texture_manifest.txt' file describing the atlas:
// Example:
//
// texture-atlas.png
// 256x256
// mine|0,64:64,64:64,128:0,128
// green ship|64,192:128,192:128,256:64,256
// ...

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

const (
	textureAtlasFile = "texture_atlas.png"
	textureManifest  = "texture_manifest.txt"
	sourceManifest   = "sources.txt"
)

func main() {

	// Open manifest file
	maniLines := loadManifest(sourceManifest)
	blocks := []*Block{}

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

			// "images/ctype/ctype0001.png"
			s := strings.Split(imagePath, "/")
			s = strings.Split(s[len(s)-1], "/")
			s = strings.Split(s[0], ".")

			name := s[0]
			iw := img.Bounds().Dx()
			ih := img.Bounds().Dy()
			fmt.Printf("(%s) %dx%d\n", name, iw, ih)

			block := &Block{
				name:  s[0],
				w:     img.Bounds().Dx(),
				h:     img.Bounds().Dy(),
				image: img,
			}
			blocks = append(blocks, block)
		}
	}

	fmt.Println("Loaded images and Blocks created")

	// This sorts by height. Even better would be to then sort by width
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[j].h < blocks[i].h
	})

	// ------------------------------------------------------
	// Now build packing tree
	// ------------------------------------------------------

	// Start small on purpose
	pow := 5.0
	packSize := int(math.Pow(2, pow)) // 5 = 32

	packer := NewPacker(packSize, packSize)

	packer.Pack(blocks)

	for !packer.Success() {
		fmt.Printf("Failed to pack into %dx%d\n", packSize, packSize)
		pow += 1.0
		if pow > 12.0 {
			fmt.Println("Stopped. Root size is limited to 4096x4096")
			break
		}
		packSize = int(math.Pow(2, pow))
		packer.Reset(packSize, packSize)
		packer.Pack(blocks)
	}

	if packer.Success() {
		fmt.Printf("Successfully packed into %dx%d\n", packSize, packSize)
		// fmt.Println(packer)
	}

	// ------------------------------------------------------
	// Now we are ready to draw the images to an atlas.
	// ------------------------------------------------------
	atlas := image.NewNRGBA(image.Rect(0, 0, packSize, packSize))

	for _, block := range blocks {
		if block.fit != nil {
			fit := block.fit
			srcImg := block.image

			target := image.Rect(fit.x, fit.y, fit.x+block.w, fit.y+block.h)

			draw.Draw(atlas, target, srcImg, srcImg.Bounds().Min, draw.Src)
		}
	}

	txf, _ := os.Create(textureAtlasFile)
	defer txf.Close()
	png.Encode(txf, atlas)

	// ------------------------------------------------------
	// Final, write texture manifest
	// ------------------------------------------------------
	f, _ := os.Create(textureManifest)
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString(textureAtlasFile + "\n")
	w.WriteString(fmt.Sprintf("%dx%d\n", packSize, packSize))

	for _, block := range blocks {
		if block.fit != nil {
			fit := block.fit
			// Write the 4 corners of the block:
			// bottom-left, bottom-right, top-right, top-left

			//       texture-space (ST "not UV")
			//     D (0,1)        (1,1) C
			//  ^      *-----------*
			//  |      |           |
			//  |+Y    |     _     |
			//  |      |           |
			//  |      |           |
			//         *-----------*
			//     A (0,0)        (1,0) B

			botlefX := fit.x // A
			botlefY := fit.y
			botrigX := fit.x + block.w // B
			botrigY := fit.y
			toprigX := fit.x + block.w // C
			toprigY := fit.y + block.h
			toplefX := fit.x // D
			toplefY := fit.y + block.h

			w.WriteString(fmt.Sprintf("%s|%d,%d:%d,%d:%d,%d:%d,%d\n",
				block.name,
				botlefX, botlefY,
				botrigX, botrigY,
				toprigX, toprigY,
				toplefX, toplefY,
			))
		}
	}
	w.Flush()
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
	packSize := int(math.Pow(2, 9)) // 8 = 256, 9 = 512 etc...

	// 128x128x4
	// 64x64x2
	// 32x32x5
	// 16x16
	blocks := []*Block{
		{w: 128, h: 128, name: "b2"},
		{w: 128, h: 128, name: "b3"},
		{w: 32, h: 32, name: "a1"},
		{w: 32, h: 32, name: "a2"},
		{w: 32, h: 32, name: "a3"},
		{w: 32, h: 32, name: "a4"},
		{w: 128, h: 128, name: "c1"},
		{w: 128, h: 128, name: "c1"},
		{w: 64, h: 64, name: "d1"},
		{w: 64, h: 64, name: "d1"},
		{w: 32, h: 32, name: "e1"},
		{w: 16, h: 16, name: "f1"},
		{w: 16, h: 16, name: "f2"},
	}

	// This sorts by height. Even better would be to then sort by width
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[j].h < blocks[i].h
	})

	packer := NewPacker(packSize, packSize)
	packer.Pack(blocks)

	var maxX, maxY int

	fmt.Println(packer)

	fmt.Printf("maxX: %d, maxY: %d\n", maxX, maxY)

	totalBlocks := len(blocks)

	fmt.Println("Container size: ", packSize)
	fmt.Println("Total blocks: ", totalBlocks)
	fmt.Println("Total packed blocks: ", packer.PackedBlockCount())
	if !packer.Success() {
		fmt.Println("##################################################")
		fmt.Println("*** Warning!! Not all blocks could be packed!! ***")
		fmt.Println("##################################################")
	} else {
		fmt.Printf("Filled: %d%%\n", packer.Efficiency())
	}

	fmt.Println("Done")
}

func testPacker2() {
	// 128x128x4
	// 64x64x2
	// 32x32x5
	// 16x16
	blocks := []*Block{
		{w: 128, h: 128, name: "b2"},
		{w: 128, h: 128, name: "b3"},
		{w: 32, h: 32, name: "a1"},
		{w: 32, h: 32, name: "a2"},
		{w: 32, h: 32, name: "a3"},
		{w: 32, h: 32, name: "a4"},
		{w: 128, h: 128, name: "c1"},
		{w: 128, h: 128, name: "c1"},
		{w: 64, h: 64, name: "d1"},
		{w: 64, h: 64, name: "d1"},
		{w: 32, h: 32, name: "e1"},
		{w: 16, h: 16, name: "f1"},
		{w: 16, h: 16, name: "f2"},
	}

	// This sorts by height. Even better would be to then sort by width
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[j].h < blocks[i].h
	})

	// Start small on purpose
	pow := 5.0
	packSize := int(math.Pow(2, pow)) // 5 = 32

	packer := NewPacker(packSize, packSize)

	packer.Pack(blocks)

	for !packer.Success() {
		fmt.Printf("Failed to pack into %dx%d\n", packSize, packSize)
		pow += 1.0
		if pow > 12.0 {
			fmt.Println("Stopped. Root size is limited to 4096x4096")
			break
		}
		packSize = int(math.Pow(2, pow))
		packer.Reset(packSize, packSize)
		packer.Pack(blocks)
	}

	if packer.Success() {
		fmt.Printf("Successfully packed into %dx%d\n", packSize, packSize)
		fmt.Println(packer)
	}
}
