package mock

import (
	"strings"
	"time"

	"github.com/megaease/easegateway/pkg/context"
	"github.com/megaease/easegateway/pkg/logger"
	"github.com/megaease/easegateway/pkg/object/httppipeline"
)

const (
	// Kind is the kind of Mock.
	Kind = "Mock"

	resultMocked = "mocked"
)

func init() {
	httppipeline.Register(&httppipeline.PluginRecord{
		Kind:            Kind,
		DefaultSpecFunc: DefaultSpec,
		NewFunc:         New,
		Results:         []string{resultMocked},
	})
}

// DefaultSpec returns default spec.
func DefaultSpec() *Spec {
	return &Spec{}
}

type (
	// Mock is plugin Mock.
	Mock struct {
		spec *Spec
		body []byte
	}

	// Spec describes the Mock.
	Spec struct {
		httppipeline.PluginMeta `yaml:",inline"`

		Rules []*Rule `yaml:"rules"`
	}

	// Rule is the mock rule.
	Rule struct {
		Path       string            `yaml:"path,omitempty" jsonschema:"omitempty,pattern=^/"`
		PathPrefix string            `yaml:"pathPrefix,omitempty" jsonschema:"omitempty,pattern=^/"`
		Code       int               `yaml:"code" jsonschema:"omitempty,format=httpcode"`
		Headers    map[string]string `yaml:"headers" jsonschema:"required"`
		Body       string            `yaml:"body" jsonschema:"omitempty"`
		Delay      string            `yaml:"delay" jsonschema:"required,format=duration"`

		delay time.Duration
	}
)

// New creates a Mock.
func New(spec *Spec, prev *Mock) *Mock {
	for _, r := range spec.Rules {
		var err error
		r.delay, err = time.ParseDuration(r.Delay)
		if err != nil {
			logger.Errorf("BUG: parse duration %s failed: %v", r.Delay, err)
		}
	}

	return &Mock{
		spec: spec,
	}
}

// Handle mocks HTTPContext.
func (m *Mock) Handle(ctx context.HTTPContext) (result string) {
	path := ctx.Request().Path()
	w := ctx.Response()

	mock := func(rule *Rule) {
		w.SetStatusCode(rule.Code)
		for key, value := range rule.Headers {
			w.Header().Set(key, value)
		}
		w.SetBody(strings.NewReader(rule.Body))
		result = resultMocked

		if rule.delay > 0 {
			logger.Debugf("delay for %v ...", rule.delay)
			time.Sleep(rule.delay)
		}
	}

	for _, rule := range m.spec.Rules {
		if rule.Path == "" && rule.PathPrefix == "" {
			mock(rule)
			return
		}

		if rule.Path == path || strings.HasPrefix(path, rule.PathPrefix) {
			mock(rule)
			return
		}
	}

	return ""
}

// Status returns status.
func (m *Mock) Status() interface{} {
	return nil
}

// Close closes Mock.
func (m *Mock) Close() {}
