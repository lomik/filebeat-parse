package parse

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

// Apr 06 15:36:08 Updated: nss-tools-3.27.1-13.el6.x86_64
// Apr 06 16:13:34 Updated: 32:bind-libs-9.8.2-0.62.rc1.el6.x86_64
// Apr 06 16:13:35 Updated: 32:bind-utils-9.8.2-0.62.rc1.el6.x86_64
// Apr 06 16:28:08 Updated: openssl-1.0.1e-57.el6.x86_64
// Apr 14 17:06:28 Installed: iotop-0.3.2-9.el6.noarch

func init() {
	Register("yum_log", ParseYumLog)
}

func ParseYumLog(text string, event common.MapStr) error {
	if len(text) < 16 {
		return errors.New("line too short")
	}

	withYear := fmt.Sprintf("%d %s", time.Now().Year(), text[:15])
	timestamp, err := time.ParseInLocation("2006 Jan 02 15:04:05", withYear, time.Local)
	if err != nil {
		return err
	}
	event["timestamp"] = timestamp.Format("2006-01-02T15:04:05.000Z0700")

	text = text[16:]

	a := strings.SplitN(text, ": ", 2)
	if len(a) != 2 {
		return ParseError
	}

	event["message"] = a[0]
	event["package"] = a[1]

	return nil
}
