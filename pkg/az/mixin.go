package az

import (
	"get.porter.sh/porter/pkg/runtime"
)

type Mixin struct {
	runtime.RuntimeConfig
	userAgent string
}

// New azure mixin client, initialized with useful defaults.
func New() *Mixin {
	return NewFor(runtime.NewConfig())
}

func NewFor(cfg runtime.RuntimeConfig) *Mixin {
	m := &Mixin{
		RuntimeConfig: cfg,
	}

	m.SetUserAgent()
	return m
}
