package az

import (
	"testing"

	"get.porter.sh/mixin/az/pkg"
	"get.porter.sh/porter/pkg/config"
	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/runtime"
)

type TestMixin struct {
	*Mixin
	TestContext *portercontext.TestContext
}

// NewTestMixin initializes a mixin test client, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	testCfg := config.NewTestConfig(t)

	// Clear this out when testing since our CI environment has modifications to it
	testCfg.Unsetenv(AzureUserAgentEnvVar)

	cfg := runtime.NewConfigFor(testCfg.Config)
	m := &TestMixin{
		Mixin:       NewFor(cfg),
		TestContext: testCfg.TestContext,
	}
	t.Cleanup(func() {
		pkg.Version = ""
		pkg.Commit = ""
	})

	return m
}
