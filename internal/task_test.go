package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimExtension(t *testing.T) {
	tests := []string{
		"a/b/file",
		"a/b/file.txt",
		"a/b/file.",
	}
	expected := []string{
		"a/b/file",
		"a/b/file",
		"a/b/file",
	}

	for i := 0; i < len(tests); i++ {
		assert.Equal(t, expected[i], trimExtension(tests[i]))
	}
}
