package az

import (
	"testing"

	"get.porter.sh/mixin/az/pkg"
	"get.porter.sh/porter/pkg/context"
)

type TestMixin struct {
	*Mixin
	TestContext *context.TestContext
}

// NewTestMixin initializes a mixin test client, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	c := context.NewTestContext(t)
	m := &TestMixin{
		Mixin: &Mixin{
			Context: c.Context,
		},
		TestContext: c,
	}
	t.Cleanup(func() {
		pkg.Version = ""
		pkg.Commit = ""
	})

	return m
}
