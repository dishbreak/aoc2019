package main

import "image"

var (
	northDir image.Point = image.Pt(0, 1)
	southDir             = image.Pt(0, -1)
	eastDir              = image.Pt(1, 0)
	westDir              = image.Pt(-1, 0)
)

type Compass struct {
	Direction   image.Point
	Left, Right *Compass
}

func NewCompass() *Compass {
	north := &Compass{
		Direction: northDir,
	}

	south := &Compass{
		Direction: southDir,
	}

	east := &Compass{
		Direction: eastDir,
	}

	west := &Compass{
		Direction: westDir,
	}

	/*
	   N
	  W E
	   S
	*/

	north.Left = west
	north.Right = east

	west.Right = north
	west.Left = south

	east.Left = north
	east.Right = south

	south.Left = east
	south.Right = west

	return north
}
