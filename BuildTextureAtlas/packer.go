package main

import (
	"fmt"
	"image"
	"math"
)

// https://codeincomplete.com/articles/bin-packing/

// Block represents an area within the root
type Block struct {
	name  string
	area  int
	w, h  int
	fit   *fittedBlock
	image image.Image
}

type fittedBlock struct {
	w, h  int
	x, y  int
	used  bool
	down  *fittedBlock
	right *fittedBlock
}

// Packer is the root
type Packer struct {
	root   *fittedBlock
	blocks []*Block
}

// NewPacker returns a new packer
func NewPacker(width, height int) *Packer {
	o := new(Packer)
	o.root = &fittedBlock{w: width, h: height}
	return o
}

// Pack attempts to place all block into root
func (p *Packer) Pack(blocks []*Block) {
	p.blocks = blocks

	for _, block := range blocks {
		block.area = block.w * block.h
		node := p.findNode(p.root, block)
		if node != nil {
			block.fit = p.splitNode(node, block)
		}
	}
}

// Reset for another attempt using new dimensions
func (p *Packer) Reset(width, height int) {
	p.root = &fittedBlock{w: width, h: height}
}

// Success checks if all blocks were packed
func (p *Packer) Success() bool {
	packedBlocks := p.PackedBlockCount()
	totalBlocks := len(p.blocks)

	return packedBlocks >= totalBlocks
}

// PackedBlockCount returns the total blocks packed
func (p *Packer) PackedBlockCount() int {
	packedBlocks := 0
	for _, block := range p.blocks {
		if block.fit != nil {
			packedBlocks++
		}
	}
	return packedBlocks
}

// Efficiency returns how much of the total area was used
func (p *Packer) Efficiency() int {
	fitTotal := 0
	for _, block := range p.blocks {
		if block.fit != nil {
			fitTotal += block.area
		}
	}

	percent := int(math.Round(100.0 * float64(fitTotal) / float64(p.root.w*p.root.h)))
	return percent
}

func (p *Packer) findNode(root *fittedBlock, block *Block) *fittedBlock {
	if root.used {
		node := p.findNode(root.right, block)
		if node == nil {
			node = p.findNode(root.down, block)
		}
		return node
	} else if (block.w <= root.w) && (block.h <= root.h) {
		return root
	}

	return nil
}

func (p *Packer) splitNode(n *fittedBlock, block *Block) *fittedBlock {
	n.used = true

	n.down = &fittedBlock{
		x: n.x, y: n.y + block.h,
		w: n.w, h: n.h - block.h,
	}

	n.right = &fittedBlock{
		x: n.x + block.w, y: n.y,
		w: n.w - block.w, h: block.h,
	}

	return n
}

func (p Packer) String() string {
	s := ""
	for _, block := range p.blocks {
		if block.fit != nil {
			fit := block.fit
			s += fmt.Sprintf(
				"'%4s' x:%03d, y: %03d, (%03dx%03d) w: %03d, h: %03d, area: %d\n",
				block.name, fit.x, fit.y, block.w, block.h, fit.w, fit.h, block.area)
		}
	}
	return s
}
