package main

// https://codeincomplete.com/articles/bin-packing/

type node struct {
	x, y  int
	w, h  int
	s, t  float64
	area  int
	fit   *node
	used  bool
	right *node
	down  *node
}

type packer struct {
	root *node
}

func (p *packer) fit(blocks []*node) {
	for _, block := range blocks {
		node := p.findNode(p.root, block.w, block.h)
		if node != nil {
			block.fit = p.splitNode(node, block.w, block.h)
		}
	}
}

func (p *packer) findNode(root *node, w, h int) *node {
	if root.used {
		node := p.findNode(root.right, w, h)
		if node == nil {
			node = p.findNode(root.down, w, h)
		}
		return node
	} else if (w <= root.w) && (h <= root.h) {
		return root
	}

	return nil
}

func (p *packer) splitNode(n *node, w, h int) *node {
	n.used = true
	n.down = &node{x: n.x, y: n.y + h, w: n.w, h: n.h - h}
	n.right = &node{x: n.x + w, y: n.y, w: n.w - w, h: h}
	n.area = w * h
	return n
}
