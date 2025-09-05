package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoopVarFunction(t *testing.T) {
	// Test the loopVar function
	usedVars := make(map[string]bool)

	// First call should return a variable name
	var1 := loopVar(usedVars)
	assert.NotEmpty(t, var1)
	assert.True(t, usedVars[var1], "Variable should be marked as used")

	// Second call should return a different variable name
	var2 := loopVar(usedVars)
	assert.NotEmpty(t, var2)
	assert.NotEqual(t, var1, var2, "Should generate different variable names")
	assert.True(t, usedVars[var2], "Second variable should be marked as used")

	// Both variables should be in the used map
	assert.Len(t, usedVars, 2)
}
