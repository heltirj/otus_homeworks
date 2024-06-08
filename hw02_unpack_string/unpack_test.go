package hw02unpackstring

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},

		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestRepeatRune(t *testing.T) {
	tests := []struct {
		name     string
		rune1    rune
		rune2    rune
		expected string
		isError  bool
	}{
		{
			name:     "a2",
			rune1:    'a',
			rune2:    '2',
			expected: "aa",
			isError:  false,
		},
		{
			name:     "у4",
			rune1:    'у',
			rune2:    '4',
			expected: "уууу",
			isError:  false,
		},
		{
			name:     "mf",
			rune1:    'm',
			rune2:    'f',
			expected: "",
			isError:  true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := repeatRune(tc.rune1, tc.rune2)
			assert.Equal(t, tc.expected, result)
			if tc.isError {
				assert.Error(t, err)
			}
		})
	}
}
