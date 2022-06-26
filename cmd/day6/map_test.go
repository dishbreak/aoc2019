package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var graphInput = []string{
	"COM)B",
	"B)C",
	"C)D",
	"D)E",
	"E)F",
	"B)G",
	"G)H",
	"D)I",
	"E)J",
	"J)K",
	"K)L",
}

func TestBuildGraph(t *testing.T) {
	graph := BuildGraph(graphInput)

	com, ok := graph["COM"]
	assert.True(t, ok)
	assert.Equal(t, "B", com.Satellites[0].Name)
}
