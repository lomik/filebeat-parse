package parse

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
)

func TestParseLTSV(t *testing.T) {
	// http://ltsv.org/

	ParseTest(t, "ltsv", `host:127.0.0.1	ident:-	user:frank	time:10/Oct/2000:13:55:36 -0700	req:GET /apache_pb.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)`,
		common.MapStr{
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
	)

	ParseTest(t, "ltsv", `empty_value:	:empty_key	:	key_only		remote_addr:127.0.0.1`,
		common.MapStr{
			"empty_value": "",
			"key_only":    "",
			"remote_addr": "127.0.0.1",
		},
	)
}

func BenchmarkParseLTSV(b *testing.B) {
	ParseBenchmark(b, "ltsv", `host:127.0.0.1	ident:-	user:frank	time:10/Oct/2000:13:55:36 -0700	req:GET /apache_pb.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)`)
}
