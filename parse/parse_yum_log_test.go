package parse

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestParseYumLog(t *testing.T) {
	table := [](struct {
		line        string
		expected    common.MapStr
		expectedErr bool
	}){
		{
			line: `Apr 06 15:36:08 Updated: 32:bind-utils-9.8.2-0.62.rc1.el6.x86_64`,
			expected: common.MapStr{
				"timestamp": time.Date(time.Now().Year(), 4, 6, 15, 36, 8, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
				"message":   "Updated",
				"package":   "32:bind-utils-9.8.2-0.62.rc1.el6.x86_64"},
		},
		{
			line: `Apr 14 17:06:28 Installed: iotop-0.3.2-9.el6.noarch`,
			expected: common.MapStr{
				"timestamp": time.Date(time.Now().Year(), 4, 14, 17, 6, 28, 0, time.Local).Format("2006-01-02T15:04:05.000Z0700"),
				"message":   "Installed",
				"package":   "iotop-0.3.2-9.el6.noarch",
			},
		},
	}

	for i := 0; i < len(table); i++ {
		m := common.MapStr{}
		err := ParseYumLog(table[i].line, m)

		if table[i].expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, table[i].expected.String(), m.String())
		}
	}
}

func BenchmarkParseYumLog(b *testing.B) {
	s := `Apr 06 15:36:08 Updated: 32:bind-utils-9.8.2-0.62.rc1.el6.x86_64`

	for i := 0; i < b.N; i++ {
		input := common.MapStr{}
		ParseYumLog(s, input)
	}
}
