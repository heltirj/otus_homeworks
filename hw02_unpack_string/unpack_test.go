package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_repeatRune(t *testing.T) {
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

func Test_createSequenceFromLetter(t *testing.T) {
	type args struct {
		letter   rune
		nextRune rune
	}
	tests := []struct {
		name       string
		args       args
		wantSeq    string
		wantOffset int
		wantErr    error
	}{
		{
			name: "next rune is a letter",
			args: args{
				letter:   'a',
				nextRune: 'b',
			},
			wantSeq:    "a",
			wantOffset: 1,
			wantErr:    nil,
		},
		{
			name: "next rune is a slash",
			args: args{
				letter:   'a',
				nextRune: '\\',
			},
			wantSeq:    "a",
			wantOffset: 1,
			wantErr:    nil,
		},
		{
			name: "next rune is a digit",
			args: args{
				letter:   'a',
				nextRune: '4',
			},
			wantSeq:    "aaaa",
			wantOffset: 2,
			wantErr:    nil,
		},
		{
			name: "next rune is invalid char",
			args: args{
				letter:   'a',
				nextRune: '_',
			},
			wantSeq:    "",
			wantOffset: 0,
			wantErr:    ErrInvalidString,
		},
		{
			name: "first rune is invalid char",
			args: args{
				letter:   '2',
				nextRune: 'b',
			},
			wantSeq:    "",
			wantOffset: 0,
			wantErr:    ErrFirstArgIsNotALetter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := createSequenceFromLetter(tt.args.letter, tt.args.nextRune)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equalf(t, tt.wantSeq, got, "createSequenceFromLetter(%v, %v)", tt.args.letter, tt.args.nextRune)
			assert.Equalf(t, tt.wantOffset, got1, "createSequenceFromLetter(%v, %v)", tt.args.letter, tt.args.nextRune)
		})
	}
}

func Test_createSequenceAfterSlash(t *testing.T) {
	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		args       args
		wantSeq    string
		wantOffset int
		wantErr    error
	}{
		{
			name: "unpack escaped digit with count",
			args: args{
				runes: []rune{'2', '2'},
			},
			wantSeq:    "22",
			wantOffset: 3,
			wantErr:    nil,
		},
		{
			name: "unpack escaped slash with count",
			args: args{
				runes: []rune{'\\', '2'},
			},
			wantSeq:    "\\\\",
			wantOffset: 3,
			wantErr:    nil,
		},
		{
			name: "unpack escaped digit without count",
			args: args{
				runes: []rune{'3'},
			},
			wantSeq:    "3",
			wantOffset: 2,
			wantErr:    nil,
		},

		{
			name: "unpack escaped slash without count",
			args: args{
				runes: []rune{'\\'},
			},
			wantSeq:    "\\",
			wantOffset: 2,
			wantErr:    nil,
		},
		{
			name: "invalid first param",
			args: args{
				runes: []rune{'n', '2'},
			},
			wantSeq:    "",
			wantOffset: 0,
			wantErr:    ErrFirstArgIsNotADigitOrSlash,
		},
		{
			name: "invalid alone param",
			args: args{
				runes: []rune{'_'},
			},
			wantSeq:    "",
			wantOffset: 0,
			wantErr:    ErrFirstArgIsNotADigitOrSlash,
		},

		{
			name: "invalid params count",
			args: args{
				runes: []rune{'a', 'b', 'c'},
			},
			wantSeq:    "",
			wantOffset: 0,
			wantErr:    ErrWrongArgsCount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := createSequenceAfterEscapeSym(tt.args.runes...)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantSeq, got)
			assert.Equal(t, tt.wantOffset, got1)
		})
	}
}
