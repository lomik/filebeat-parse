package parse

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestParseNginxErrorLog(t *testing.T) {
	// some sample rows from https://gadelkareem.com/2012/07/01/nginx-error-log-reader/ example.log

	table := [](struct {
		line        string
		expected    common.MapStr
		expectedErr bool
	}){
		{
			line: `2012/04/15 22:01:47 [error] 3258#0: *887 upstream timed out (110: Connection timed out) while reading upstream, client: 192.168.126.1, server: *.example.com, request: "GET /wp-admin/options-general.php?page=ozh_ta&action=import_all&time=1334482870&_wpnonce=89590fa285 HTTP/1.1", upstream: "fastcgi://unix:/var/run/php-fpm/php-fpm.sock:", host: "www.example.com", referrer: "http://www.example.com/wp-login.php?redirect_to=http%3A%2F%2Fwww.example.com%2Fwp-admin%2Foptions-general.php%3Fpage%3Dozh_ta%26action%3Dimport_all%26time%3D1334482870%26_wpnonce%3D89590fa285&reauth=1"`,
			expected: common.MapStr{
				"timestamp": time.Date(2012, 4, 15, 22, 1, 47, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
				"level":     "error",
				"pid":       3258,
				"tid":       0,
				"sid":       887,
				"message":   "upstream timed out (110: Connection timed out) while reading upstream",
				"client":    "192.168.126.1",
				"server":    "*.example.com",
				"request":   "GET /wp-admin/options-general.php?page=ozh_ta&action=import_all&time=1334482870&_wpnonce=89590fa285 HTTP/1.1",
				"upstream":  "fastcgi://unix:/var/run/php-fpm/php-fpm.sock:",
				"http_host": "www.example.com",
				"referrer":  "http://www.example.com/wp-login.php?redirect_to=http%3A%2F%2Fwww.example.com%2Fwp-admin%2Foptions-general.php%3Fpage%3Dozh_ta%26action%3Dimport_all%26time%3D1334482870%26_wpnonce%3D89590fa285&reauth=1",
			},
		},
		{
			line: `2017/05/16 16:53:25 [notice] 26016#0: using inherited sockets from "6;"`,
			expected: common.MapStr{
				"timestamp": time.Date(2017, 5, 16, 16, 53, 25, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
				"level":     "notice",
				"pid":       26016,
				"tid":       0,
				"message":   "using inherited sockets from \"6;\"",
			},
		},
	}

	for i := 0; i < len(table); i++ {
		m := common.MapStr{}
		err := ParseNginxErrorLog(table[i].line, m)

		if table[i].expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, table[i].expected.String(), m.String())
		}
	}
}

func BenchmarkParseNginxErrorLog(b *testing.B) {
	s := `2012/04/15 22:01:47 [error] 3258#0: *887 upstream timed out (110: Connection timed out) while reading upstream, client: 192.168.126.1, server: *.example.com, request: "GET /wp-admin/options-general.php?page=ozh_ta&action=import_all&time=1334482870&_wpnonce=89590fa285 HTTP/1.1", upstream: "fastcgi://unix:/var/run/php-fpm/php-fpm.sock:", host: "www.example.com", referrer: "http://www.example.com/wp-login.php?redirect_to=http%3A%2F%2Fwww.example.com%2Fwp-admin%2Foptions-general.php%3Fpage%3Dozh_ta%26action%3Dimport_all%26time%3D1334482870%26_wpnonce%3D89590fa285&reauth=1"`

	for i := 0; i < b.N; i++ {
		input := common.MapStr{}
		ParseNginxErrorLog(s, input)
	}
}
