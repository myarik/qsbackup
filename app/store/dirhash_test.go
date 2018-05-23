package store

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestDirHash(t *testing.T) {
	hash1a, err := DirHash("../../test/testdata/hash1")
	require.NoError(t, err)
	hash1b, err := DirHash("../../test/testdata/hash1")
	require.NoError(t, err)
	hash2a, err := DirHash("../../test/testdata/hash2")
	require.NoError(t, err)
	require.Equal(t, hash1a, hash1b, "hash1 and hash1b should be identical")
	require.NotEqual(t, hash1a, hash2a, "hash1 and hash2 should not be the same")
}


func TestDirID(t *testing.T) {
	hash1a := DirID("/home/ubuntu/test")
	require.Equal(t, hash1a, "1319342438", "hash1a is invalid")
	hash1b := DirID("/home/ubuntu/test2")
	require.NotEqual(t, hash1a, hash1b, "hash1 and hash2 should not be the same")
}
