package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputStream := make(chan int64, 1) // must be buffered because last input never gets read.
	output := make(chan int64)
	defer close(inputStream)

	result := runRobot(ctx, inputStream, output, 0)

	go func() {
		defer close(output)
		i := <-inputStream
		assert.Equal(t, int64(0), i)
		output <- 1
		output <- 0

		i = <-inputStream
		assert.Equal(t, int64(0), i)
		output <- 0
		output <- 0

		<-inputStream
		output <- 1
		output <- 0

		<-inputStream
		output <- 1
		output <- 0

		<-inputStream
		output <- 0
		output <- 1

		<-inputStream
		output <- 1
		output <- 0

		<-inputStream
		output <- 1
		output <- 0

	}()

	assert.Equal(t, 6, len(<-result))
}
