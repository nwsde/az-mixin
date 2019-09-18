package az

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestMixin_Build(t *testing.T) {
		const buildOutput = `RUN apt-get update && apt-get install -y apt-transport-https lsb-release gnupg curl
RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.asc.gpg
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list
RUN apt-get update && apt-get install -y azure-cli
`

	t.Run("build with config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")

		wantOutput := buildOutput + `RUN az extension add -y --name iot
`
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build without config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-without-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")

		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, buildOutput, gotOutput)
	})
}