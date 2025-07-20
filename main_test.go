package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AlwaysPass(t *testing.T) {
	assert.True(t, true)
}
