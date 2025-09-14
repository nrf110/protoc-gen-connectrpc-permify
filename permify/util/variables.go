package util

import (
	"fmt"
	"sync/atomic"
)

// Global counter for deterministic variable name generation
var variableCounter atomic.Uint64

// ResetVariableCounter resets the counter - useful for tests or when starting a new file
func ResetVariableCounter() {
	variableCounter.Store(0)
}

// VariableName generates a deterministic variable name using a counter
// The names follow the pattern: v1, v2, v3, etc.
func VariableName() string {
	counter := variableCounter.Add(1)
	return fmt.Sprintf("v%d", counter)
}
