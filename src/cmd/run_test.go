package cmd

import (
	"github.com/rocktavious/autopilot/v2022"
	"testing"
)

func TestExample(t *testing.T) {
	// Arrange
	data := `{}`
	// Act

	// Assert
	autopilot.Equals(t, `{}`, data)
}
