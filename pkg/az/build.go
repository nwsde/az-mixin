package az

import (
	"context"
	"fmt"
	"text/template"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v3"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the az mixin in porter.yaml
// mixins:
// - az:
//     extensions:
//     - NAME
type MixinConfig struct {
	// UserAgentOptOut allows a bundle author to opt out from adding porter and the mixin's version to the az CLI user agent string.
	UserAgentOptOut bool `yaml:"userAgentOptOut,omitempty"`

	// ClientVersion is the version of the az CLI to install.
	ClientVersion string `yaml:"clientVersion,omitempty"`

	// Extensions is a list of az CLI extensions to install.
	Extensions []string `yaml:"extensions,omitempty"`
}

// buildConfig is the set of configuration options for the mixin's portion of the Dockerfile
type buildConfig struct {
	MixinConfig

	// AzureUserAgent is the contents of the az CLI user agent environment variable.
	AzureUserAgent string
}

// The package version of the az cli follows this format:
// VERSION-1~DISTRO_CODENAME So if we are running on debian stretch and have a
// version of 1.2.3, the package version would be 1.2.3-1~stretch.
const buildTemplate string = `
ENV PORTER_AZ_MIXIN_USER_AGENT_OPT_OUT="{{ .UserAgentOptOut}}"
ENV AZURE_HTTP_USER_AGENT="{{ .AzureUserAgent }}"
RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
	apt-get update && apt-get install -y apt-transport-https lsb-release gnupg curl
RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.asc.gpg
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list
RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
	apt-get update && apt-get install -y --no-install-recommends \
	{{ if eq .ClientVersion ""}}azure-cli{{else}}azure-cli={{.ClientVersion}}-1~$(lsb_release -cs){{end}}
{{ range $ext := .Extensions }}
RUN az extension add -y --name {{ $ext }}
{{ end }}
`

// Build installs the az cli and any configured extensions.
func (m *Mixin) Build(ctx context.Context) error {
	var input BuildInput
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	tmpl, err := template.New("dockerfile").Parse(buildTemplate)
	if err != nil {
		return fmt.Errorf("error parsing Dockerfile template for the az mixin: %w", err)
	}

	cfg := buildConfig{MixinConfig: input.Config}
	if !input.Config.UserAgentOptOut {
		// If they opt out, then the user agent string defaulted in the bundle will be empty
		cfg.AzureUserAgent = m.userAgent
	}

	if err = tmpl.Execute(m.Out, cfg); err != nil {
		return fmt.Errorf("error generating Dockerfile lines for the az mixin: %w", err)
	}

	return nil
}
