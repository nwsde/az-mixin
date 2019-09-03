package az

import "fmt"

// Build installs the az cli
func (m *Mixin) Build() error {
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y apt-transport-https lsb-release gnupg curl`)
	fmt.Fprintln(m.Out, `RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.asc.gpg`)
	fmt.Fprintln(m.Out, `RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list`)
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y azure-cli`)
	return nil
}
