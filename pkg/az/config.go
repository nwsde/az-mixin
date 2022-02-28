package az

import (
	"strings"

	"get.porter.sh/porter/pkg"
)

const AZURE_HTTP_USER_AGENT = "AZURE_HTTP_USER_AGENT"

func (m *Mixin) SetUserAgent() {
	value := []string{pkg.UserAgent(), m.UserAgent()}

	if agentStr, ok := m.LookupEnv(AZURE_HTTP_USER_AGENT); ok {
		value = append(value, agentStr)
	}

	m.Setenv(AZURE_HTTP_USER_AGENT, strings.Join(value, " "))
}

func (m *Mixin) UserAgent() string {
	v := m.Version()
	return v.Name + "/" + v.Version
}
