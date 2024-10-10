package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	wantEnv := Environment{
		"BAR": EnvValue{
			Value: "bar",
		},
		"EMPTY": EnvValue{},
		"FOO":   EnvValue{Value: "   foo\nwith new line"},
		"HELLO": EnvValue{Value: "\"hello\""},
		"UNSET": EnvValue{},
	}

	require.Equal(t, 0, RunCmd([]string{"testdata/echo.sh"}, wantEnv))
}
