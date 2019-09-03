package az

import (
	"fmt"

	"github.com/deislabs/porter-az/pkg"
)

func (m *Mixin) PrintVersion() {
	fmt.Fprintf(m.Out, "az mixin %s (%s)\n", pkg.Version, pkg.Commit)
}
