package parse

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

func TestParseYumLog(t *testing.T) {
	ParseTest(t, "yum_log",
		`Apr 06 15:36:08 Updated: 32:bind-utils-9.8.2-0.62.rc1.el6.x86_64`,
		common.MapStr{
			"timestamp": time.Date(time.Now().Year(), 4, 6, 15, 36, 8, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
			"message":   "Updated",
			"package":   "32:bind-utils-9.8.2-0.62.rc1.el6.x86_64",
		},
	)

	ParseTest(t, "yum_log",
		`Apr 14 17:06:28 Installed: iotop-0.3.2-9.el6.noarch`,
		common.MapStr{
			"timestamp": time.Date(time.Now().Year(), 4, 14, 17, 6, 28, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
			"message":   "Installed",
			"package":   "iotop-0.3.2-9.el6.noarch",
		},
	)
}

func BenchmarkParseYumLog(b *testing.B) {
	ParseBenchmark(b, "yum_log", `Apr 06 15:36:08 Updated: 32:bind-utils-9.8.2-0.62.rc1.el6.x86_64`)
}
