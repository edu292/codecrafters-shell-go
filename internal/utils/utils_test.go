package utils

import (
	"fmt"
	"testing"
)

func TestExpandPath(t *testing.T) {
	absPath := ExpandPath("../../app/")
	fmt.Printf("t: %v\n", absPath)
}
