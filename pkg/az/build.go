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
	ClientVersion string   `yaml:"clientVersion,omitempty"`
	Extensions    []string `yaml:"extensions,omitempty"`
}

// The package version of the az cli follows this format:
// VERSION-1~DISTRO_CODENAME So if we are running on debian stretch and have a
// version of 1.2.3, the package version would be 1.2.3-1~stretch.
const buildTemplate string = `
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

	if err = tmpl.Execute(m.Out, input.Config); err != nil {
		return fmt.Errorf("error generating Dockerfile lines for the az mixin: %w", err)
	}

	return nil
}
