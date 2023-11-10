package main

import (
	"regexp"
	"strconv"
)

var pointExp regexp.Regexp = *regexp.MustCompile(`<x=(-*\d+), y=(-*\d+), z=(-*\d+)>`)

func FromString(input string) (p [3]int) {
	matches := pointExp.FindStringSubmatch(input)
	if matches == nil {
		return
	}
	for i := range p {
		p[i], _ = strconv.Atoi(matches[i+1])
	}
	return
}
