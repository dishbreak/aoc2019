package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompass(t *testing.T) {
	t.Run("full clockwise rotation", func(t *testing.T) {
		c := NewCompass()
		result := c.Right
		result = result.Right
		result = result.Right
		result = result.Right

		assert.Equal(t, c.Direction, result.Direction)
	})
	t.Run("full counterclockwise rotation", func(t *testing.T) {
		c := NewCompass()
		result := c.Left
		result = result.Left
		result = result.Left
		result = result.Left

		assert.Equal(t, c.Direction, result.Direction)
	})
}
