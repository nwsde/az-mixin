package az

import (
	"strconv"
	"strings"

	"get.porter.sh/porter/pkg"
)

const (
	// AzureUserAgentEnvVar is the environment variable used by the az CLI to set
	// the user agent string sent to Azure.
	AzureUserAgentEnvVar = "AZURE_HTTP_USER_AGENT"

	// UserAgentOptOutEnvVar is the name of the environment variable that disables
	// user agent reporting.
	UserAgentOptOutEnvVar = "PORTER_AZ_MIXIN_USER_AGENT_OPT_OUT"
)

// SetUserAgent sets the AZURE_HTTP_USER_AGENT environment variable with
// the full user agent string, which includes both a portion for porter and the
// mixin.
func (m *Mixin) SetUserAgent() {
	// Check if PORTER_AZ_MIXIN_USER_AGENT_OPT_OUT=true, which disables editing the user agent string
	if optOut, _ := strconv.ParseBool(m.Getenv(UserAgentOptOutEnvVar)); optOut {
		return
	}

	// Check if we have already set the user agent
	if m.userAgent != "" {
		return
	}

	// Append porter and the mixin's version to the user agent string. Some clouds and
	// environments will have set the environment variable already and we don't want
	// to clobber it.
	porterUserAgent := pkg.UserAgent()
	value := []string{porterUserAgent, m.GetMixinUserAgent()}
	if agentStr, ok := m.LookupEnv(AzureUserAgentEnvVar); ok {
		value = append(value, agentStr)
	}

	m.userAgent = strings.Join(value, " ")

	// Set the az CLI user agent as an environment variable so that when we call the
	// az CLI, it's automatically passed too.
	m.Setenv(AzureUserAgentEnvVar, m.userAgent)
}

// GetMixinUserAgent returns the portion of the user agent string for the mixin.
func (m *Mixin) GetMixinUserAgent() string {
	v := m.Version()
	return "getporter/" + v.Name + "/" + v.Version
}
