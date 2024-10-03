package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

var testDir = "testdata/env/"

func TestReadDir(t *testing.T) {
	t.Run("starts with empty envs", func(t *testing.T) {
		env := Environment{
			"BAR": EnvValue{
				Value: "bar",
			},
			"EMPTY": EnvValue{},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{},
		}

		for k := range env {
			err := os.Unsetenv(k)
			if err != nil {
				t.Fatal(err)
			}
		}

		res, err := ReadDir(testDir)
		require.NoError(t, err)
		require.Equal(t, env, res)
	})
	t.Run("starts with set envs", func(t *testing.T) {
		env := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: true,
			},
			"EMPTY": EnvValue{
				NeedRemove: true,
			},
			"FOO": EnvValue{Value: "   foo\nwith new line",
				NeedRemove: true,
			},
			"HELLO": EnvValue{Value: "\"hello\"",
				NeedRemove: true,
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
		}

		for k := range env {
			err := os.Setenv(k, "test")
			if err != nil {
				t.Fatal(err)
			}
		}

		res, err := ReadDir(testDir)
		require.NoError(t, err)
		require.Equal(t, env, res)
	})

	t.Run("invalid file name", func(t *testing.T) {
		file, err := os.Create(path.Join(testDir, "test=test"))
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(file.Name())
		defer file.Close()

		_, err = ReadDir(testDir)
		require.ErrorIs(t, err, ErrInvalidFileName)
	})
}
