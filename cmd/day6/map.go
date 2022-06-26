package main

import (
	"strings"
)

type Planet struct {
	Name       string
	Orbits     int
	Parent     *Planet
	Satellites []*Planet
	Visited    bool
	Steps      int
}

func upsert(graph map[string]*Planet, name string) *Planet {
	p, ok := graph[name]
	if ok {
		return p
	}

	p = &Planet{
		Name:       name,
		Satellites: make([]*Planet, 0),
	}

	graph[name] = p

	return p
}

func link(parent, child *Planet) {
	child.Parent = parent
	parent.Satellites = append(parent.Satellites, child)
}

func BuildGraph(input []string) map[string]*Planet {
	result := make(map[string]*Planet)

	for _, line := range input {
		parts := strings.Split(line, ")")
		parent := upsert(result, parts[0])
		child := upsert(result, parts[1])
		link(parent, child)
	}

	return result
}
