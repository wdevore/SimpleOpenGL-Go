package textures

import (
	"bufio"
	"image"
	"image/draw"
	"os"
	"strconv"
	"strings"
)

// TextureCoord is:
// uv coordinates start from the upper left corner (v-axis is facing down).
// st coordinates start from the lower left corner (t-axis is facing up).
// s = u;
// t = 1-v;
type TextureCoord struct {
	S, T float32
}

// SubTexture is a block within the image atlas
type SubTexture struct {
	name          string
	textureCoords []*TextureCoord
}

// NewSubTexture creates a
func NewSubTexture(name string) *SubTexture {
	o := new(SubTexture)
	o.name = name
	o.textureCoords = []*TextureCoord{}
	return o
}

// TextureAtlas contains an image atlas
type TextureAtlas struct {
	manifest      string
	width, height int64
	atlas         *image.NRGBA

	subTextures []*SubTexture
}

// NewTextureAtlas creates a new atlas
func NewTextureAtlas(manifest string) *TextureAtlas {
	o := new(TextureAtlas)
	o.manifest = manifest
	o.subTextures = []*SubTexture{}

	return o
}

// Build setups the atlas based on manifest
func (t *TextureAtlas) Build() {
	manifestFile, err := os.Open(t.manifest)
	if err != nil {
		panic(err)
	}

	defer manifestFile.Close()

	lines := []string{}

	scanner := bufio.NewScanner(manifestFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	textureFile := lines[0]

	t.atlas, err = t.loadImage(textureFile)
	if err != nil {
		panic(err)
	}

	s := strings.Split(lines[1], "x")
	t.width, _ = strconv.ParseInt(s[0], 10, 64)
	t.height, _ = strconv.ParseInt(s[1], 10, 64)

	for i := 2; i < len(lines); i++ {
		s = strings.Split(lines[i], "|")

		ts := NewSubTexture(s[0])

		coords := strings.Split(s[1], ":")

		for _, coord := range coords {
			xy := strings.Split(coord, ",")
			x, _ := strconv.ParseInt(xy[0], 10, 64)
			y, _ := strconv.ParseInt(xy[1], 10, 64)

			sc := float32(x) / float32(t.width)
			tc := float32(y) / float32(t.height)
			// fmt.Println(sc, ",", tc)
			ts.textureCoords = append(ts.textureCoords, &TextureCoord{S: sc, T: tc})
		}

		t.subTextures = append(t.subTextures, ts)
	}
}

// Atlas returns image atlas
func (t *TextureAtlas) Atlas() *image.NRGBA {
	return t.atlas
}

// TextureCoords returns the assigned coords of named sub texture
func (t *TextureAtlas) TextureCoords(name string) []*TextureCoord {
	for _, subTex := range t.subTextures {
		if name == subTex.name {
			return subTex.textureCoords
		}
	}

	return nil
}

func (t *TextureAtlas) loadImage(path string) (*image.NRGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	r := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	flippedImg := image.NewNRGBA(r)

	// Transfer data to image
	draw.Draw(nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src)

	// Flip horizontally or around Y-axis
	// for j := 0; j < nrgba.Bounds().Dy(); j++ {
	// 	for i := 0; i < nrgba.Bounds().Dx(); i++ {
	// 		flippedImg.Set(bounds.Dx()-i, j, nrgba.At(i, j))
	// 	}
	// }

	// Flip vertically or around the X-axis
	for j := 0; j < nrgba.Bounds().Dy(); j++ {
		for i := 0; i < nrgba.Bounds().Dx(); i++ {
			flippedImg.Set(i, bounds.Dy()-j, nrgba.At(i, j))
		}
	}

	return flippedImg, nil
}
