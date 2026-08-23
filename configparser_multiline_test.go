package configparser_test

import (
	"bytes"
	"testing"

	"github.com/bigkevmcd/go-configparser"
)

// BenchmarkParseLargeMultilineValue measures parsing an option with a large
// multiline value. Each continuation line used to be appended with string
// concatenation, which reallocates and copies the whole accumulated value on
// every line, making parsing O(n^2) in the number of continuation lines.
func BenchmarkParseLargeMultilineValue(b *testing.B) {
	const n = 300000

	var buf bytes.Buffer
	buf.WriteString("[section]\nkey = first\n")
	for i := 0; i < n; i++ {
		buf.WriteString(" cont\n")
	}
	input := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := configparser.ParseReader(bytes.NewReader(input)); err != nil {
			b.Fatalf("ParseReader returned error: %v", err)
		}
	}
}
