package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFile_IsEmpty(t *testing.T) {
	empty := File{}
	require.True(t, empty.IsEmpty())

	notEmpty := File{
		ID:   "id",
		Name: "not-empty.txt",
		Path: "path1/path2",
		User: 1,
		Size: 15,
	}
	require.False(t, notEmpty.IsEmpty())
}
