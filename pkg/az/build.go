package az

import (
	"context"
	"fmt"

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
	Extensions []string
}

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

	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y apt-transport-https lsb-release gnupg curl`)
	fmt.Fprintln(m.Out, `RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.asc.gpg`)
	fmt.Fprintln(m.Out, `RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list`)
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y azure-cli`)

	for _, ext := range input.Config.Extensions {
		fmt.Fprintf(m.Out, "RUN az extension add -y --name %s\n", ext)
	}

	return nil
}
