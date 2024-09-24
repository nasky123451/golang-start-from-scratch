package main

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type Block struct {
	id         uuid.UUID
	start, end time.Time
}

func (b Block) Duration() time.Duration {
	return b.end.Sub(b.start)
}

type Group struct {
	Block

	children []uuid.UUID
}

func (g *Group) Merge(b Block) {
	if g.end.IsZero() || g.end.Before(b.end) {
		g.end = b.end
	}
	if g.start.IsZero() || g.start.After(b.start) {
		g.start = b.start
	}
	g.children = append(g.children, b.id)
}

func Compact(blocks ...Block) Block {
	sort.Sort(sortable(blocks)) // Sort the blocks

	g := &Group{}
	g.id = uuid.New()
	for _, b := range blocks {
		g.Merge(b)
	}
	return g.Block
}

type sortable []Block

func (s sortable) Len() int {
	return len(s)
}

func (s sortable) Less(i, j int) bool {
	// Sort by the start time of the blocks
	return s[i].start.Before(s[j].start)
}

func (s sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}