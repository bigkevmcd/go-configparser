package configparser_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bigkevmcd/go-configparser"
)

// TestParseLargeMultilineValue guards against quadratic-time behaviour when a
// single option has a large multiline value. Each continuation line used to be
// appended to the value with string concatenation, which reallocates and copies
// the whole accumulated value on every line, making parsing O(n^2) in the number
// of continuation lines. A ~600KB input took tens of seconds to parse.
func TestParseLargeMultilineValue(t *testing.T) {
	const n = 300000

	var b bytes.Buffer
	b.WriteString("[section]\nkey = first\n")
	for i := 0; i < n; i++ {
		b.WriteString(" cont\n")
	}

	start := time.Now()
	p, err := configparser.ParseReader(bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatalf("ParseReader returned error: %v", err)
	}
	elapsed := time.Since(start)

	value, err := p.Get("section", "key")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	// "first" plus n continuation lines => n newlines in the joined value.
	if got := strings.Count(value, "\n"); got != n {
		t.Fatalf("value has %d newlines, want %d", got, n)
	}

	// Linear parsing completes in well under a second; the previous quadratic
	// implementation needed roughly a minute for this input. The generous bound
	// only fails on quadratic (or worse) behaviour.
	if elapsed > 10*time.Second {
		t.Fatalf("parsing %d continuation lines took %v, expected near-linear time", n, elapsed)
	}
}
