package parse

import (
	"strings"

	"github.com/elastic/beats/libbeat/common"
)

func init() {
	Register("ltsv", ParseLTSV)
}

func ParseLTSV(text string, event common.MapStr) error {
	arr := strings.Split(text, "\t")
	for i := 0; i < len(arr); i++ {
		if len(arr[i]) == 0 {
			continue
		}

		p1 := strings.IndexByte(arr[i], ':')
		if p1 < 0 {
			event[arr[i]] = ""
			continue
		}
		// skipe empty key
		if p1 == 0 {
			continue
		}
		if p1 == len(arr[i])-1 {
			event[arr[i][:p1]] = ""
			continue
		}

		event[arr[i][:p1]] = arr[i][p1+1:]
	}

	return nil
}
