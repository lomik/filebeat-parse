package parse

import (
	"strings"
	"testing"

	"github.com/elastic/beats/libbeat/common"
)

func ParseTest(t *testing.T, parser string, line string, expected common.MapStr) {
	h := parser
	p1 := strings.IndexByte(parser, ',')
	if p1 >= 0 {
		h = parser[:p1]
	}
	handler := ParseHandler[h]

	m := common.MapStr{}
	err := handler(line, m)

	if expected == nil {
		if err == nil {
			t.Fatal("error not raised")
		}
	} else {
		if err != nil {
			t.Fatal(err)
		}

		if m.String() != expected.String() {
			t.Fatalf("%s (actual) != %s (expected)", m.String(), expected.String())
		}
	}
}

func ParseBenchmark(b *testing.B, parser string, line string) {
	h := parser
	p1 := strings.IndexByte(parser, ',')
	if p1 >= 0 {
		h = parser[:p1]
	}
	handler := ParseHandler[h]

	for i := 0; i < b.N; i++ {
		m := common.MapStr{}
		handler(line, m)
	}
}
