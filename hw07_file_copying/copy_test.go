package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var destPath = path.Join(os.TempDir(), "test_dst")

func TestCopy(t *testing.T) {
	t.Run("empty source", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.ErrorIs(t, err, ErrUndefinedSourceFile)
	})

	t.Run("empty dest", func(t *testing.T) {
		f, n, err := CreateTestFile()
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		err = Copy(f.Name(), destPath, n+10, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		f, n, err := CreateTestFile()
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		err = Copy(f.Name(), destPath, n+10, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("limit exceeds file size", func(t *testing.T) {
		f, n, err := CreateTestFile()
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		err = Copy(f.Name(), destPath, 0, n+10)
		require.Nil(t, err)

		destStat, err := os.Stat(destPath)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, destStat.Size(), n)
	})
}

func CreateTestFile() (*os.File, int64, error) {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		return nil, 0, err
	}

	n, err := f.WriteString("В чёрном-чёрном городе сидели чёрные-чёрные люди")
	if err != nil {
		return nil, 0, err
	}

	return f, int64(n), nil
}
