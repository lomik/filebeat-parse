package parse

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestParseLTSV(t *testing.T) {
	// http://ltsv.org/

	table := [](struct {
		line        string
		expected    common.MapStr
		expectedErr bool
	}){
		{
			line: `host:127.0.0.1	ident:-	user:frank	time:10/Oct/2000:13:55:36 -0700	req:GET /apache_pb.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)`,
			expected: common.MapStr{
				"host":    "127.0.0.1",
				"ident":   "-",
				"user":    "frank",
				"time":    "10/Oct/2000:13:55:36 -0700",
				"req":     "GET /apache_pb.gif HTTP/1.0",
				"status":  "200",
				"size":    "2326",
				"referer": "http://www.example.com/start.html",
				"ua":      "Mozilla/4.08 [en] (Win98; I ;Nav)",
			},
		},
		{
			line: `empty_value:	:empty_key	:	key_only		remote_addr:127.0.0.1`,
			expected: common.MapStr{
				"empty_value": "",
				"key_only":    "",
				"remote_addr": "127.0.0.1",
			},
		},
	}

	for i := 0; i < len(table); i++ {
		m := common.MapStr{}
		err := ParseLTSV(table[i].line, m)

		if table[i].expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, table[i].expected.String(), m.String())
		}
	}
}

func BenchmarkParseLTSV(b *testing.B) {
	s := `host:127.0.0.1	ident:-	user:frank	time:10/Oct/2000:13:55:36 -0700	req:GET /apache_pb.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)`

	for i := 0; i < b.N; i++ {
		input := common.MapStr{}
		ParseLTSV(s, input)
	}
}
