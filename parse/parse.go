package parse

import (
	"fmt"
	"strings"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
	"github.com/pkg/errors"
)

var ParseHandler = make(map[string]func(text string, event common.MapStr) error)

func Register(name string, callback func(text string, event common.MapStr) error) {
	ParseHandler[name] = callback
}

type Parse struct {
	Field    string
	handlers map[string]func(text string, event common.MapStr) error
	mu       *sync.RWMutex
}

type ParseConfig struct {
	Field string `config:"field"`
}

var (
	ParseDefaultConfig = ParseConfig{
		Field: "fields.parse",
	}
)

func init() {
	processors.RegisterPlugin("parse",
		configChecked(NewParse,
			allowedFields("field", "when")))
}

var DropEvent = errors.New("DropEvent")
var ParseError = errors.New("parse error")

func NewParse(c *common.Config) (processors.Processor, error) {
	config := ParseDefaultConfig

	err := c.Unpack(&config)

	if err != nil {
		logp.Warn("Error unpacking config for parse")
		return nil, fmt.Errorf("fail to unpack the parse configuration: %s", err)
	}

	f := Parse{
		Field:    config.Field,
		handlers: make(map[string]func(text string, event common.MapStr) error),
		mu:       new(sync.RWMutex),
	}
	return f, nil
}

func (f Parse) Run(event *beat.Event) (*beat.Event, error) {
	var errs []string

	parserObj, err := event.GetValue(f.Field)
	if err != nil {
		if errors.Cause(err) == common.ErrKeyNotFound {
			return event, nil
		}
		errs = append(errs, err.Error())
		return event, fmt.Errorf(strings.Join(errs, ", "))
	}

	parser, ok := parserObj.(string)
	if !ok {
		errs = append(errs, fmt.Sprintf("wrong parser value: %#v", parserObj))
		return event, fmt.Errorf(strings.Join(errs, ", "))
	}

	data, err := event.GetValue("message")

	if err != nil {
		if errors.Cause(err) == common.ErrKeyNotFound {
			return event, nil
		}
		errs = append(errs, err.Error())
		return event, fmt.Errorf(strings.Join(errs, ", "))
	}

	text, ok := data.(string)

	event.Delete(f.Field)
	event.Delete("message")

	var handler func(text string, event common.MapStr) error

	f.mu.RLock()
	handler, ok = f.handlers[parser]
	f.mu.RUnlock()

	if !ok {
		f.mu.Lock()
		handler, ok = f.handlers[parser]
		if !ok {
			class := parser
			p1 := strings.IndexByte(parser, ',')
			if p1 >= 0 {
				class = parser[:p1]
			}
			handler = ParseHandler[class]

			if handler == nil {
				errs = append(errs, fmt.Sprintf("bad parser %#v", parser))
			} else {
				f.handlers[parser] = handler
			}
		}
		f.mu.Unlock()

		if len(errs) > 0 {
			return event, fmt.Errorf(strings.Join(errs, ", "))
		}
	}

	err = handler(text, event.Fields)

	if err != nil {
		if err == DropEvent {
			return nil, nil
		}

		event.Fields["parse_error"] = err.Error()
		event.Fields["message"] = text
	}

	if len(errs) > 0 {
		return event, fmt.Errorf(strings.Join(errs, ", "))
	}
	return event, nil
}

func (f Parse) String() string {
	return "parse=" + f.Field
}
