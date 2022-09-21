package az

import (
	"get.porter.sh/porter/pkg/runtime"
)

type Mixin struct {
	runtime.RuntimeConfig
}

// New azure mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	m := &Mixin{
		RuntimeConfig: runtime.NewConfig(),
	}

	m.SetUserAgent()
	return m, nil

}
