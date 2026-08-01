package cupwriter

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestWriterFlush(t *testing.T) {
	testCases := []struct {
		name  string
		lines int
	}{
		{"single", 1},
		{"double", 2},
		{"multi", 3},
	}

	out := new(bytes.Buffer)
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			defer out.Reset()
			input := strings.Repeat("foo\n", test.lines)
			cw := New(out, true)
			// need at least 2 flushes for esc chars to appear
			for range 2 {
				cw.WriteString(input)
				err := cw.Flush(test.lines)
				if err != nil {
					t.Fatal(err)
				}
			}
			p, err := parse(out.String())
			if err != nil {
				t.Fatal(err)
			}
			if p.input != input {
				t.Errorf("want: %q, got %q", input, p.input)
			}
			if p.lines != test.lines {
				t.Errorf("want: %d, got %d", test.lines, p.lines)
			}
		})
	}
}

type parsed struct {
	input string
	lines int
}

func parse(escaped string) (p *parsed, err error) {
	p = new(parsed)
	for i, sep := range []string{escOpen, cuuAndEd} {
		if before, after, found := strings.Cut(escaped, sep); found {
			switch i {
			case 0:
				p.input, escaped = before, after
			case 1:
				p.lines, err = strconv.Atoi(before)
			}
		}
	}
	return
}
