package az

import (
	"testing"

	"get.porter.sh/mixin/az/pkg"
	"github.com/stretchr/testify/require"
)

func TestSetUserAgent(t *testing.T) {
	pkg.Commit = "abc123"
	pkg.Version = "v1.2.3"

	m := NewTestMixin(t)
	m.SetUserAgent()

	expected := "porter az/" + pkg.Version
	require.Contains(t, m.Getenv(AZURE_HTTP_USER_AGENT), expected)
}
