package az

import (
	"testing"

	"get.porter.sh/mixin/az/pkg"
	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/runtime"
)

type TestMixin struct {
	*Mixin
	TestContext *portercontext.TestContext
}

// NewTestMixin initializes a mixin test client, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	c := portercontext.NewTestContext(t)

	// Clear this out when testing since our CI environment has modifications to it
	c.Unsetenv(AZURE_HTTP_USER_AGENT)

	cfg := runtime.NewConfigFor(c.Context)
	m := &TestMixin{
		Mixin:       NewFor(cfg),
		TestContext: c,
	}
	t.Cleanup(func() {
		pkg.Version = ""
		pkg.Commit = ""
	})

	return m
}
